package http

import (
	"ap4y.me/cloud-mail/config"
	"ap4y.me/cloud-mail/notmuch"
)

type AccountData struct {
	Address   string         `json:"address"`
	Name      string         `json:"name"`
	Mailboxes []MailboxStats `json:"mailboxes"`
	Tags      []string       `json:"tags"`
}

type MailboxStats struct {
	config.Mailbox
	Unread int `json:"unread"`
}

type Threads struct {
	HasMore bool             `json:"has_more"`
	Threads []notmuch.Thread `json:"threads"`
}
