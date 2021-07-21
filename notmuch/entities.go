package notmuch

type Thread struct {
	ID           string   `json:"thread"`
	Timestamp    int64    `json:"timestamp"`
	DateRelative string   `json:"date_relative"`
	Matched      int      `json:"matched"`
	Total        int      `json:"total"`
	Authors      string   `json:"authors"`
	Subject      string   `json:"subject"`
	Query        []string `json:"query"`
	Tags         []string `json:"tags"`
}

type MessageThread []MessageTuple

type MessageTuple struct {
	Message
	MessageThread
}

type Message struct {
	ID           string                 `json:"id"`
	Match        bool                   `json:"match"`
	Excluded     bool                   `json:"excluded"`
	Filename     []string               `json:"filename"`
	Timestamp    int64                  `json:"timestamp"`
	DateRelative string                 `json:"date_relative"`
	Tags         []string               `json:"tags"`
	Body         []MessageBody          `json:"body"`
	Headers      map[string]interface{} `json:"headers"`
}

type MessageBody struct {
	ID                 int         `json:"id"`
	ContentType        string      `json:"content-type"`
	ContentDisposition string      `json:"content-disposition"`
	Content            interface{} `json:"content"`
}

type Reply struct {
	ReplyHeaders map[string]interface{} `json:"reply-headers"`
	Original     Message                `json:"original"`
}
