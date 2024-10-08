package http

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"ap4y.me/amail/config"
	"ap4y.me/amail/notmuch"
	"ap4y.me/amail/smtp"
	"ap4y.me/amail/tagger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func SetLogger(l zerolog.Logger) {
	logger = l
}

type Server struct {
	http.Server

	client     *notmuch.Client
	smtpClient *smtp.Client
	refresher  tagger.Refresher

	name      string
	addresses []string
	mailboxes []config.Mailbox
	tags      []string
}

func NewServer(
	name string, addresses []string, mailboxes []config.Mailbox, tags []string,
	smtpClient *smtp.Client, refresher tagger.Refresher, staticBundle fs.FS,
	creds map[string]string,
) (*Server, error) {

	c, err := notmuch.NewClient()
	if err != nil {
		return nil, err
	}

	s := &Server{
		client: c, smtpClient: smtpClient, name: name,
		addresses: addresses, mailboxes: mailboxes,
		refresher: refresher, tags: tags,
	}

	r := chi.NewRouter()
	if len(creds) != 0 {
		r.Use(middleware.BasicAuth("mail", creds))
	}

	r.Route("/api", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				logger.Debug().Msgf("%s %s", r.Method, r.URL.Path)
				next.ServeHTTP(w, r)
			})
		})

		r.Get("/mailboxes", s.mailboxesHandler)
		r.Get("/search/{term}", s.searchHandler)
		r.Get("/address/{term}", s.addressHandler)
		r.Put("/tags", s.tagsHandler)
		r.Get("/threads/{threadID}", s.threadHandler)

		r.Post("/messages", s.sendMessageHandler)
		r.Post("/drafts", s.saveMessageHandler)
		r.Get("/messages/{messageID}/reply", s.messageReplyHandler)
		r.Get("/messages/{messageID}/parts/{partID}", s.messagePartsHandler)
		r.Get("/messages/{messageID}/w3m/{partID}", s.messageW3mHandler)
	})

	fs := http.FileServer(http.FS(staticBundle))
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
	data := AccountData{s.addresses[0], s.name, make([]MailboxStats, len(s.mailboxes)), s.tags}

	if err := s.refresher.RefreshMailboxes(); err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	for idx, mailbox := range s.mailboxes {
		if mailbox.Tracked {
			unread, err := s.client.Count(mailbox.Terms+" and tag:unread", notmuch.CountOutputMessages)
			if err != nil {
				sendError(w, r, err, http.StatusBadRequest)
				return
			}

			data.Mailboxes[idx] = MailboxStats{mailbox, unread}
		} else {
			data.Mailboxes[idx] = MailboxStats{mailbox, 0}
		}
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		sendError(w, r, err, http.StatusBadRequest)
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

	threads, err := s.client.Search(term, perPage+1, page*perPage)
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	res := Threads{len(threads) == perPage+1, threads}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		sendError(w, r, err, http.StatusBadRequest)
	}
}

func (s *Server) addressHandler(w http.ResponseWriter, r *http.Request) {
	term := chi.URLParam(r, "term")
	addresses, err := s.client.Address(term)
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(addresses); err != nil {
		sendError(w, r, err, http.StatusBadRequest)
	}
}

func (s *Server) threadHandler(w http.ResponseWriter, r *http.Request) {
	threadID := chi.URLParam(r, "threadID")

	messages, err := s.client.Show("thread:" + threadID)
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	if len(messages) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(messages[0]); err != nil {
		sendError(w, r, err, http.StatusBadRequest)
	}
}

func (s *Server) messagePartsHandler(w http.ResponseWriter, r *http.Request) {
	base64ID := chi.URLParam(r, "messageID")
	partID := chi.URLParam(r, "partID")

	messageID, err := base64.StdEncoding.DecodeString(base64ID)
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	attachment, part, err := s.client.Attachment(string(messageID), partID)
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	if cType, ok := part["content-type"].(string); ok {
		if part["content-charset"] != nil {
			cType += "; charset=" + part["content-charset"].(string)
		}
		w.Header().Set("Content-Type", cType)
	}

	filename := "attachment"
	if part["filename"] != nil {
		if name, ok := part["filename"].(string); ok {
			filename = name
		} else if name, ok := part["filename"].([]string); ok {
			filename = name[0]
		}
	}

	http.ServeContent(w, r, filename, time.Now(), attachment)
}

func (s *Server) messageW3mHandler(w http.ResponseWriter, r *http.Request) {
	base64ID := chi.URLParam(r, "messageID")
	partID := chi.URLParam(r, "partID")

	messageID, err := base64.StdEncoding.DecodeString(base64ID)
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	attachment, part, err := s.client.Attachment(string(messageID), partID)
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	cmd := exec.Command("w3m", "-T", part["content-type"].(string), "-dump")
	cmd.Stdin = attachment
	out, err := cmd.Output()
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"text": string(out)}); err != nil {
		sendError(w, r, err, http.StatusBadRequest)
	}
}

func (s *Server) messageReplyHandler(w http.ResponseWriter, r *http.Request) {
	base64ID := chi.URLParam(r, "messageID")

	messageID, err := base64.StdEncoding.DecodeString(base64ID)
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	replyTo := r.URL.Query().Get("reply-to")
	reply, err := s.client.Reply("id:"+string(messageID), notmuch.ReplyToType(replyTo))
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	reply.ReplyHeaders["From"] = fmt.Sprintf("%s <%s>", s.name, s.addresses[0])
	if err := json.NewEncoder(w).Encode(reply); err != nil {
		sendError(w, r, err, http.StatusBadRequest)
	}
}

func (s *Server) tagsHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Terms string   `json:"terms"`
		Tags  []string `json:"tags"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := s.client.Tag(data.Terms, data.Tags); err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	newTags, err := s.client.Dump(data.Terms)
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(newTags); err != nil {
		sendError(w, r, err, http.StatusBadRequest)
	}
}

func (s *Server) sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	msg, err := s.messageFromForm(r)
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	defer func(attached []*smtp.Attachment) {
		for _, attach := range attached {
			if err := attach.Close(); err != nil {
				logger.Debug().Msgf("Failed to close attachment: %v", err)
			}
		}
	}(msg.Attachments)

	m, err := s.smtpClient.Send(msg)
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	var mbox *config.Mailbox
	for _, mailbox := range s.mailboxes {
		if mailbox.ID == "sent" {
			mbox = &mailbox
			break
		}
	}

	if mbox != nil && mbox.Folder != "" {
		if err := s.client.Insert(mbox.Folder, m, "+sent", "-inbox", "-unread"); err != nil {
			sendError(w, r, err, http.StatusBadRequest)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
		sendError(w, r, err, http.StatusBadRequest)
	}
}

func (s *Server) saveMessageHandler(w http.ResponseWriter, r *http.Request) {
	msg, err := s.messageFromForm(r)
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	defer func(attached []*smtp.Attachment) {
		for _, attach := range attached {
			if err := attach.Close(); err != nil {
				logger.Debug().Msgf("Failed to close attachment: %v", err)
			}
		}
	}(msg.Attachments)

	m, _, _, err := s.smtpClient.Compose(msg)
	if err != nil {
		sendError(w, r, err, http.StatusBadRequest)
		return
	}

	var mbox *config.Mailbox
	for _, mailbox := range s.mailboxes {
		if mailbox.ID == "draft" {
			mbox = &mailbox
			break
		}
	}

	if mbox != nil && mbox.Folder != "" {
		if err := s.client.Insert(mbox.Folder, m, "+draft", "-inbox", "-unread"); err != nil {
			sendError(w, r, err, http.StatusBadRequest)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(map[string]string{}); err != nil {
		sendError(w, r, err, http.StatusBadRequest)
	}
}

func (s *Server) messageFromForm(r *http.Request) (*smtp.Message, error) {
	if err := r.ParseMultipartForm(1024 * 1024); err != nil {
		return nil, err
	}

	form := r.MultipartForm
	attachments := make([]*smtp.Attachment, 0)
	for _, attach := range form.File["attachments[]"] {
		f, err := attach.Open()
		if err != nil {
			return nil, err
		}

		attachments = append(attachments, &smtp.Attachment{
			ReadCloser:  f,
			Filename:    attach.Filename,
			ContentType: attach.Header.Get("Content-Type"),
		})
	}

	for _, attach := range form.Value["attachments[]"] {
		items := strings.Split(attach, ":")
		if len(items) != 2 {
			return nil, errors.New("invalid attachment: id is not valid")
		}

		attachment, part, err := s.client.Attachment(items[0], items[1])
		if err != nil {
			return nil, err
		}

		attachments = append(attachments, &smtp.Attachment{
			ReadCloser:  io.NopCloser(attachment),
			Filename:    part["filename"].(string),
			ContentType: part["content-type"].(string),
		})
	}

	return &smtp.Message{
		To:          form.Value["to[]"],
		CC:          form.Value["cc[]"],
		Body:        r.FormValue("body"),
		Subject:     r.FormValue("subject"),
		Attachments: attachments,
		Headers:     formMap(form, "headers"),
	}, nil
}

func sendError(w http.ResponseWriter, r *http.Request, err error, status int) {
	logger.Info().Err(err).Msgf("%s %s: %d", r.Method, r.URL.Path, http.StatusBadRequest)
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func formMap(form *multipart.Form, fKey string) map[string]string {
	headers := map[string]string{}

	for key, val := range form.Value {
		if strings.HasPrefix(key, fKey+"[") {
			headerKey := strings.ReplaceAll(
				strings.ReplaceAll(key, fKey+"[", ""), "]", "",
			)
			headers[headerKey] = val[0]
		}
	}

	return headers
}
