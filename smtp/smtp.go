package smtp

import (
	"fmt"

	"ap4y.me/cloud-mail/config"
	"github.com/rs/zerolog"
	"gopkg.in/gomail.v2"
)

var logger zerolog.Logger

func SetLogger(l zerolog.Logger) {
	logger = l
}

type Authenticator interface {
	Password(username, hostname string) (string, error)
}

type Message struct {
	To      []string `json:"to"`
	CC      []string `json:"cc"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

type Client struct {
	address  string
	username string
	hostname string
	port     int

	auth Authenticator
}

func New(address string, conf config.Submission, auth Authenticator) *Client {
	return &Client{address, conf.Username, conf.Hostname, conf.Port, auth}
}

func (c *Client) Send(msg *Message) error {
	logger.Debug().Msgf("sending %#v", msg)

	m := gomail.NewMessage()
	m.SetHeader("From", c.address)
	m.SetHeader("To", msg.To...)
	if len(msg.CC) > 0 {
		m.SetHeader("Cc", msg.CC...)
	}
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/plain", msg.Body)

	pass, err := c.auth.Password(c.username, c.hostname)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	d := gomail.NewDialer(c.hostname, c.port, c.username, pass)
	return d.DialAndSend(m)
}
