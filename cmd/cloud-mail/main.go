package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"ap4y.me/cloud-mail/config"
	"ap4y.me/cloud-mail/http"
	"ap4y.me/cloud-mail/smtp"
	"ap4y.me/cloud-mail/static/public"
	"ap4y.me/cloud-mail/tagger"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
)

const notmuchConfigPath = ".notmuch-config"

var (
	logger = zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log    = logger.With().Str("sys", "main").Timestamp().Logger()

	confPath = flag.String("config", "config.toml", "path to the config file")
	logLevel = flag.String("log", "debug", "log level for logger")
)

func main() {
	flag.Parse()

	lvl, err := zerolog.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal().Msgf("Failed to parse log level: %s", err)
	}

	log = log.Level(lvl)
	tagger.SetLogger(logger.Level(lvl).With().Str("sys", "tag").Timestamp().Logger())
	http.SetLogger(logger.Level(lvl).With().Str("sys", "http").Timestamp().Logger())
	smtp.SetLogger(logger.Level(lvl).With().Str("sys", "smtp").Timestamp().Logger())

	conf, err := config.FromFile(*confPath)
	if err != nil {
		log.Fatal().Msgf("Failed to parse config: %s", err)
	}

	ensureNotmuchConfig(conf)

	t, err := tagger.New(conf.TagRules, conf.Cleanup.Tags)
	if err != nil {
		log.Fatal().Msgf("Error creating a tagger: %s", err)
	}
	if err := t.RefreshMailboxes(); err != nil {
		log.Fatal().Msgf("Failed to refresh mailboxes: %s", err)
	}

	setupRefresh(conf, t)
	setupCleanup(conf.Cleanup.Interval.Duration, t)

	if len(conf.Addresses) == 0 {
		log.Fatal().Msg("Specify at least one address")
	}

	client := smtp.New(
		fmt.Sprintf("%s <%s>", conf.Name, conf.Addresses[0]),
		conf.Submission,
		&conf.Submission,
	)

	s, err := http.NewServer(conf.Name, conf.Addresses, conf.Mailboxes, client, public.Content)
	if err != nil {
		log.Fatal().Msgf("Error creating an http server: %s", err)
	}

	s.Addr = ":8000"
	log.Info().Msg("Starting on: " + s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal().Msgf("Startup error: %s", err)
	}
}

func setupRefresh(conf *config.Config, t *tagger.Tagger) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal().Msgf("Failed to create FS watcher: %s", err)
	}
	defer watcher.Close()
	for _, folder := range conf.Refresh.Watch {
		path := filepath.Join(conf.Maildir, folder)
		if err := watcher.Add(path); err != nil {
			log.Fatal().Msgf("Failed to add %s to FS watcher: %s", path, err)
		}
		log.Info().Msgf("Added %s to FS watcher", path)
	}

	go func() {
		for event := range watcher.Events {
			log.Info().Msgf("FS event in folder %s, event: %d", event.Name, event.Op)
			if err := t.RefreshMailboxes(); err != nil {
				log.Warn().Err(err).Msg("Failed to refresh mailboxes")
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(conf.Refresh.Interval.Duration)
		for range ticker.C {
			if err := t.RefreshMailboxes(); err != nil {
				log.Warn().Err(err).Msg("Failed to refresh mailboxes")
			}
		}
	}()
}

func setupCleanup(interval time.Duration, t *tagger.Tagger) {
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			if err := t.Cleanup(); err != nil {
				log.Warn().Err(err).Msg("Failed to cleanup mailboxes")
			}
		}
	}()
}

func ensureNotmuchConfig(conf *config.Config) {
	path := os.Getenv("NOTMUCH_CONFIG")
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Msgf("Failed to resolve home dir path: %s", err)
		}

		path = filepath.Join(homeDir, notmuchConfigPath)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return
	}

	log.Info().Msgf("%s does not exist. Will generate a new one.", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatal().Msgf("Failed to create %s: %s", path, err)
	}

	if err := conf.WriteNotmuchConfig(f); err != nil {
		log.Fatal().Msgf("Failed to write %s: %s", path, err)
	}
}
