package tagger

import (
	"strings"

	"ap4y.me/cloud-mail/notmuch"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func SetLogger(l zerolog.Logger) {
	logger = l
}

type Tagger struct {
	client   *notmuch.Client
	tagRules map[string]string
}

func New(tagRules map[string]string) (*Tagger, error) {
	c, err := notmuch.NewClient()
	if err != nil {
		return nil, err
	}

	return &Tagger{c, tagRules}, nil
}

func (t *Tagger) RefreshMailboxes() error {
	logger.Debug().Msg("Re-indexing maildir")
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
