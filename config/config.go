package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

var ErrInvalidPasswordCommand = errors.New("invalid PasswordCommand")

type Config struct {
	Name       string
	Addresses  []string
	Maildir    string
	Mailboxes  []Mailbox
	TagRules   map[string]string `toml:"tag_rules"`
	Submission Submission
	Refresh    Refresh
	Cleanup    Cleanup
}

func FromFile(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	var conf Config
	if _, err := toml.Decode(string(data), &conf); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return &conf, nil
}

type Mailbox struct {
	ID         string   `json:"id"`
	Folder     string   `json:"folder"`
	Terms      string   `json:"terms"`
	ToggleTags []string `json:"tags" toml:"toggle_tags"`
}

type Submission struct {
	Hostname string
	Port     int

	Username        string
	PasswordCommand string `toml:"password_command"`
}

func (s *Submission) Password(username, hostname string) (string, error) {
	if s.PasswordCommand == "" {
		return "", ErrInvalidPasswordCommand
	}

	fields := strings.Fields(s.PasswordCommand)
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

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type Refresh struct {
	Watch    []string
	Interval Duration
}

type Cleanup struct {
	Tags     []string
	Interval Duration
}
