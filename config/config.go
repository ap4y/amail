package config

import (
	"errors"
	"os/exec"
	"strings"
	"time"
)

var ErrInvalidPasswordCommand = errors.New("invalid PasswordCommand")

type Config struct {
	Name            string
	Addresses       []string
	Maildir         string
	Mailboxes       []Mailbox
	TagRules        map[string]string
	RefreshInterval time.Duration
	Submission      Submission
	PasswordCommand string
}

func (c *Config) Password(username, hostname string) (string, error) {
	if c.PasswordCommand == "" {
		return "", ErrInvalidPasswordCommand
	}

	fields := strings.Fields(c.PasswordCommand)
	if len(fields) == 0 {
		return "", ErrInvalidPasswordCommand
	}

	cmd := exec.Command(fields[0], fields[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

type Mailbox struct {
	ID         string   `json:"id"`
	Folder     string   `json:"folder"`
	Terms      string   `json:"terms"`
	ToggleTags []string `json:"tags"`
}

type Submission struct {
	Hostname string
	Port     int

	Username string
}
