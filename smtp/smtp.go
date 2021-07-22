package smtp

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"ap4y.me/cloud-mail/config"
	"github.com/rs/zerolog"
	"gopkg.in/gomail.v2"
)

var logger zerolog.Logger
var userAgent = "amail 0.0.1"

func SetLogger(l zerolog.Logger) {
	logger = l
}

func SetUserAgent(ua string) {
	userAgent = ua
}

type Authenticator interface {
	Password(username, hostname string) (string, error)
}

type Headers map[string]string

type Message struct {
	To      []string `json:"to"`
	CC      []string `json:"cc"`
	Headers Headers  `json:"headers"`
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

func (c *Client) Send(msg *Message) (*gomail.Message, error) {
	logger.Debug().Msgf("sending %#v", msg)

	m := gomail.NewMessage()
	m.SetHeader("User-agent", userAgent)
	m.SetHeader("From", c.address)
	m.SetHeader("To", msg.To...)
	if len(msg.CC) > 0 {
		m.SetHeader("Cc", msg.CC...)
	}
	m.SetHeader("Subject", msg.Subject)
	if msg.Headers != nil {
		for key, val := range msg.Headers {
			m.SetHeader(key, val)
		}
	}

	m.SetHeader("Message-ID", c.generateMessageId())
	m.SetBody("text/plain", msg.Body)

	pass, err := c.auth.Password(c.username, c.hostname)
	if err != nil {
		return nil, fmt.Errorf("auth: %w", err)
	}

	d := gomail.NewDialer(c.hostname, c.port, c.username, pass)
	return m, d.DialAndSend(m)
}

func (c *Client) generateMessageId() string {
	items := strings.Split(c.address, "@")

	id := make([]byte, 5)
	rand.Read(id)

	return fmt.Sprintf("<%s.amail@%s>", hex.EncodeToString(id), items[1])
}
