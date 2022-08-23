package tagger

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"ap4y.me/amail/notmuch"
	"github.com/emersion/go-msgauth/dkim"
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
	logger.Debug().Msg("Re-indexing maildir")
	for {
		err := t.client.Index()
		if err == nil {
			break
		}

		var e *exec.ExitError
		if !errors.As(err, &e) || e.ExitCode() != 75 {
			break
		}
	}

	t.processDkimTags()

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

		files, err := t.client.SearchWithOutput("tag:"+tag, "files")
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

func (t *Tagger) processDkimTags() {
	messages, err := t.client.SearchWithOutput(
		"tag:unread and tag:inbox not tag:dkim-none not tag:dkim-fail not tag:dkim-ok",
		"messages",
	)
	if err != nil {
		logger.Info().Err(err).Msg("DKIM tagging failed")
		return
	}

	for _, messageId := range messages {
		if err := t.addDkimTags(messageId); err != nil {
			logger.Info().Err(err).Msgf("DKIM tagging failed: %s", messageId)
		}
	}
}

func (t *Tagger) addDkimTags(messageId string) error {
	logger.Debug().Msgf("DKIM tagging %s", messageId)

	r, _, err := t.client.Attachment(messageId, "0")
	if err != nil {
		return fmt.Errorf("open raw: %v", err)
	}

	verifications, err := dkim.Verify(r)
	if err != nil {
		return fmt.Errorf("dkim: %v", err)
	}

	tag := "dkim-ok"
	if len(verifications) == 0 {
		tag = "dkim-none"
	} else {
		for _, v := range verifications {
			if v.Err != nil {
				tag = "dkim-fail"
				break
			}
		}
	}

	if err := t.client.Tag("id:"+messageId, []string{"+" + tag}); err != nil {
		return fmt.Errorf("tag: %v", err)
	}

	return nil
}
