package http

import "ap4y.me/cloud-mail/notmuch"

type Mailbox struct {
	ID    string   `json:"id"`
	Terms string   `json:"terms"`
	Tags  []string `json:"tags"`
}

type AccountData struct {
	Address   string         `json:"address"`
	Mailboxes []MailboxStats `json:"mailboxes"`
}

type MailboxStats struct {
	Mailbox
	Unread int `json:"unread"`
}

type Threads struct {
	Total   int              `json:"total"`
	Threads []notmuch.Thread `json:"threads"`
}
