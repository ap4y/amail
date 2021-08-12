package config

import "time"

type Config struct {
	Name            string
	Addresses       []string
	Maildir         string
	Mailboxes       []Mailbox
	TagRules        map[string]string
	RefreshInterval time.Duration
	Submission      Submission
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
