package tagger

import (
	"os"
	"strings"

	"ap4y.me/cloud-mail/notmuch"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func SetLogger(l zerolog.Logger) {
	logger = l
}

type Refresher interface {
	RefreshMailboxes() error
}

type Tagger struct {
	client      *notmuch.Client
	tagRules    map[string]string
	cleanupTags []string
}

func New(tagRules map[string]string, cleanupTags []string) (*Tagger, error) {
	c, err := notmuch.NewClient()
	if err != nil {
		return nil, err
	}

	return &Tagger{c, tagRules, cleanupTags}, nil
}

func (t *Tagger) RefreshMailboxes() error {
	logger.Info().Msg("Re-indexing maildir")
	if err := t.client.Index(); err != nil {
		return err
	}

	for tags, terms := range t.tagRules {
		if err := t.client.Tag(terms, strings.Fields(tags)); err != nil {
			logger.Info().Err(err).Msgf("Failed to tag %s with %s", terms, tags)
		}
	}

	return nil
}

func (t *Tagger) Cleanup() error {
	for _, tag := range t.cleanupTags {
		logger.Debug().Msgf("Performing cleanup for tag %s", tag)

		files, err := t.client.SearchFiles("tag:" + tag)
		if err != nil {
			logger.Info().Err(err).Msgf("Failed to cleanup tag %s: %s", tag, err)
			continue
		}

		for _, file := range files {
			logger.Debug().Msgf("Removing file %s", file)
			if err := os.Remove(file); err != nil {
				logger.Info().Err(err).Msgf("Failed to remove file %s: %s", file, err)
				continue
			}
		}
	}

	return t.RefreshMailboxes()
}
