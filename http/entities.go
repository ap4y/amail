package http

import "ap4y.me/cloud-mail/notmuch"

type Mailboxes struct {
	Address   string                  `json:"address"`
	Mailboxes map[string]MailboxStats `json:"mailboxes"`
}

type MailboxStats struct {
	ID     string `json:"id"`
	Terms  string `json:"terms"`
	Unread int    `json:"unread"`
	Total  int    `json:"total"`
}

type Threads struct {
	Total   int              `json:"total"`
	Threads []notmuch.Thread `json:"threads"`
}
