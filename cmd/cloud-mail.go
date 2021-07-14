package main

import (
	"os"
	"time"

	"ap4y.me/cloud-mail/config"
	"ap4y.me/cloud-mail/http"
	"ap4y.me/cloud-mail/tagger"
	"github.com/rs/zerolog"
)

var conf = config.Config{
	Addresses: []string{"mail@ap4y.me", "ap4y@me.com", "arthur.evstifeev@gmail.com"},
	Mailboxes: []config.Mailbox{
		{"inbox", "tag:personal and tag:inbox", []string{"inbox"}},
		{"archive", "tag:archive", []string{"archive"}},
		{"sent", "tag:sent", []string{"sent"}},
		{"spam", "tag:spam", []string{"spam"}},
		{"trash", "tag:trash", []string{"trash"}},
		{"openbsd", "to:tech@openbsd.org and tag:inbox", []string{"inbox"}},
		{"unknown", "not tag:personal and tag:inbox", []string{"inbox"}},
	},
	TagRules: map[string]string{
		"+personal":               "to:mail@ap4y.me or to:ap4y@me.com",
		"+archive -unread -inbox": "folder:Archive",
		"+sent -unread -inbox":    "folder:Sent",
		"+spam -unread -inbox":    "folder:Junk",
		"+trash -unread -inbox":   "folder:Trash",
	},
	RefreshInterval: time.Minute,
}

var logger = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr})

func main() {
	log := logger.With().Str("sys", "main").Timestamp().Logger()
	tagger.SetLogger(logger.With().Str("sys", "tag").Timestamp().Logger())
	http.SetLogger(logger.With().Str("sys", "http").Timestamp().Logger())

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

	s, err := http.NewServer(conf.Addresses, conf.Mailboxes)
	if err != nil {
		log.Fatal().Msgf("Error creating an http server: %s", err)
	}

	s.Addr = ":8000"
	log.Info().Msg("Starting on: " + s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal().Msgf("Startup error: %s", err)
	}
}
