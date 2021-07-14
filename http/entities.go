package http

import (
	"ap4y.me/cloud-mail/config"
	"ap4y.me/cloud-mail/notmuch"
)

type AccountData struct {
	Address   string         `json:"address"`
	Mailboxes []MailboxStats `json:"mailboxes"`
}

type MailboxStats struct {
	config.Mailbox
	Unread int `json:"unread"`
}

type Threads struct {
	Total   int              `json:"total"`
	Threads []notmuch.Thread `json:"threads"`
}
