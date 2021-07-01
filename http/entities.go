package http

type Mailboxes struct {
	Address   string                  `json:"address"`
	Mailboxes map[string]MailboxStats `json:"mailboxes"`
}

type MailboxStats struct {
	ID     string `json:"id"`
	Folder string `json:"folder"`
	Unread int    `json:"unread"`
	Total  int    `json:"total"`
}
