package config

import "time"

type Config struct {
	Addresses       []string
	Mailboxes       []Mailbox
	TagRules        map[string]string
	RefreshInterval time.Duration
}

type Mailbox struct {
	ID         string   `json:"id"`
	Terms      string   `json:"terms"`
	ToggleTags []string `json:"tags"`
}
