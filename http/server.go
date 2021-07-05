package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"ap4y.me/cloud-mail/notmuch"
	"github.com/go-chi/chi"
)

var mailboxes = map[string]string{
	"inbox":   "INBOX",
	"archive": "Archive",
	"sent":    "Sent",
	"drafts":  "Drafts",
	"junk":    "Junk",
	"trash":   "Trash",
}

type Server struct {
	http.Server

	client       *notmuch.Client
	primaryEmail string
}

func NewServer(primaryEmail string) (*Server, error) {
	c, err := notmuch.NewClient()
	if err != nil {
		return nil, err
	}

	s := &Server{client: c, primaryEmail: primaryEmail}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/mailboxes", s.mailboxesHandler)
		r.Get("/search/{term}", s.searchHandler)
		r.Get("/threads/{threadID}", s.threadHandler)
		r.Get("/messages/{messageID}/parts/{partID}", s.messagePartsHandler)
	})

	fs := http.FileServer(http.Dir("./static/public")) // TODO: replace with embed
	for mailbox := range mailboxes {
		r.Handle(fmt.Sprintf("/%s*", mailbox), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = "/"
			fs.ServeHTTP(w, r)
		}))
	}
	r.Handle("/search*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/"
		fs.ServeHTTP(w, r)
	}))
	r.Handle("/*", fs)

	s.Handler = r
	return s, nil
}

func (s *Server) mailboxesHandler(w http.ResponseWriter, r *http.Request) {
	data := Mailboxes{s.primaryEmail, map[string]MailboxStats{}}

	for mailbox, folder := range mailboxes {
		unread, err := s.client.Count("folder:" + folder + " tag:unread")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		total, err := s.client.Count("folder:" + folder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data.Mailboxes[mailbox] = MailboxStats{mailbox, "folder:" + folder, unread, total}
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	term := chi.URLParam(r, "term")
	params := r.URL.Query()

	perPage := 50
	if val, err := strconv.Atoi(params.Get("per")); err == nil {
		perPage = val
	}

	page := 0
	if val, err := strconv.Atoi(params.Get("page")); err == nil {
		page = val
	}

	threads, err := s.client.Search(term, perPage, page*perPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	total, err := s.client.Count(term)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := Threads{total, threads}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) threadHandler(w http.ResponseWriter, r *http.Request) {
	threadID := chi.URLParam(r, "threadID")

	messages, err := s.client.Show("thread:" + threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(messages[0][0]); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Server) messagePartsHandler(w http.ResponseWriter, r *http.Request) {
	base64ID := chi.URLParam(r, "messageID")
	partID := chi.URLParam(r, "partID")

	messageID, err := base64.StdEncoding.DecodeString(base64ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	attachment, err := s.client.Attachment(string(messageID), partID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.ServeContent(w, r, "attachment", time.Now(), attachment)
}
