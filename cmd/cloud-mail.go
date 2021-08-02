package main

import (
	"os"
	"time"

	"ap4y.me/cloud-mail/config"
	"ap4y.me/cloud-mail/http"
	"ap4y.me/cloud-mail/smtp"
	"ap4y.me/cloud-mail/tagger"
	"github.com/rs/zerolog"
)

var conf = config.Config{
	Addresses: []string{"mail@ap4y.me", "ap4y@me.com", "arthur.evstifeev@gmail.com"},
	Mailboxes: []config.Mailbox{
		{"inbox", "INBOX", "tag:personal and tag:inbox", []string{"inbox"}},
		{"archive", "Archive", "tag:archive", []string{"archive"}},
		{"sent", "Sent", "tag:sent", []string{"sent"}},
		{"draft", "Drafts", "tag:draft", []string{"draft"}},
		{"spam", "Junk", "tag:spam", []string{"spam"}},
		{"trash", "Trash", "tag:trash", []string{"trash"}},
		{"openbsd", "", "tag:openbsd and tag:inbox", []string{"inbox"}},
		{"unsorted", "", "not tag:personal and not tag:list and tag:inbox", []string{"inbox"}},
	},
	TagRules: map[string]string{
		"+personal":               "to:mail@ap4y.me or to:ap4y@me.com",
		"+openbsd +list":          "to:tech@openbsd.org",
		"+archive -unread -inbox": "folder:Archive",
		"+sent -unread -inbox":    "folder:Sent",
		"+spam -unread -inbox":    "folder:Junk",
		"+trash -unread -inbox":   "folder:Trash",
	},
	RefreshInterval: time.Minute,
	Submission: config.Submission{
		Hostname: "mail.ap4y.me",
		Port:     587,
		Username: "ap4y",
	},
}

var logger = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr})

func main() {
	log := logger.With().Str("sys", "main").Timestamp().Logger()
	tagger.SetLogger(logger.With().Str("sys", "tag").Timestamp().Logger())
	http.SetLogger(logger.With().Str("sys", "http").Timestamp().Logger())
	smtp.SetLogger(logger.With().Str("sys", "smtp").Timestamp().Logger())

	t, err := tagger.New(conf.TagRules)
	if err != nil {
		log.Fatal().Msgf("Error creating a tagger: %s", err)
	}
	if err := t.RefreshMailboxes(); err != nil {
		log.Fatal().Msgf("Failed to refresh mailboxes: %s", err)
	}

	go func() {
		ticker := time.NewTicker(conf.RefreshInterval)
		for range ticker.C {
			if err := t.RefreshMailboxes(); err != nil {
				log.Warn().Err(err).Msg("Failed to refresh mailboxes")
			}
		}
	}()

	if len(conf.Addresses) == 0 {
		log.Fatal().Msg("Specify at least one address")
	}

	client := smtp.New(conf.Addresses[0], conf.Submission, auth(func(username, hostname string) (string, error) {
		return "crews96/gust", nil
	}))

	s, err := http.NewServer(conf.Addresses, conf.Mailboxes, client)
	if err != nil {
		log.Fatal().Msgf("Error creating an http server: %s", err)
	}

	s.Addr = ":8000"
	log.Info().Msg("Starting on: " + s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal().Msgf("Startup error: %s", err)
	}
}

type auth func(string, string) (string, error)

func (a auth) Password(username, hostname string) (string, error) {
	return a(username, hostname)
}
