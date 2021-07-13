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

var mailboxes = []Mailbox{
	{"inbox", "tag:inbox to:mail@ap4y.me"},
	{"archive", "tag:archive"},
	{"sent", "tag:sent"},
	{"spam", "tag:spam"},
	{"trash", "tag:trash"},
	{"openbsd", "to:tech@openbsd.org and tag:inbox"},
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
		r.Put("/messages/{messageID}/tags", s.messageTagsHandler)
	})

	fs := http.FileServer(http.Dir("./static/public")) // TODO: replace with embed
	for _, mailbox := range mailboxes {
		r.Handle(fmt.Sprintf("/%s*", mailbox.ID), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	data := AccountData{s.primaryEmail, make([]MailboxStats, len(mailboxes))}

	for idx, mailbox := range mailboxes {
		unread, err := s.client.Count(mailbox.Terms + " and tag:unread")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data.Mailboxes[idx] = MailboxStats{mailbox, unread}
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

	if err := json.NewEncoder(w).Encode(messages[0]); err != nil {
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

func (s *Server) messageTagsHandler(w http.ResponseWriter, r *http.Request) {
	base64ID := chi.URLParam(r, "messageID")

	messageID, err := base64.StdEncoding.DecodeString(base64ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var tags []string
	if err := json.NewDecoder(r.Body).Decode(&tags); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.client.Tag("id:"+string(messageID), tags); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTags, err := s.client.Dump("id:" + string(messageID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(newTags); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func messageTags(message []interface{}) (tags []interface{}) {
	tags = make([]interface{}, 0)

	thread, ok := message[0].([]interface{})
	if !ok {
		return
	}

	threadMessage, ok := thread[0].(map[string]interface{})
	if !ok {
		return
	}

	return threadMessage["tags"].([]interface{})
}
