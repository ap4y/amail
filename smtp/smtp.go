package smtp

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"ap4y.me/cloud-mail/config"
	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog"
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

type Attachment struct {
	io.ReadCloser
	Filename string
}

type Message struct {
	To          []string      `json:"to"`
	CC          []string      `json:"cc"`
	Headers     Headers       `json:"headers"`
	Subject     string        `json:"subject"`
	Body        string        `json:"body"`
	Attachments []*Attachment `json:"-"`
}

type Client struct {
	address  *mail.Address
	username string
	hostname string
	port     int

	auth Authenticator
}

func New(address *mail.Address, conf config.Submission, auth Authenticator) *Client {
	return &Client{address, conf.Username, conf.Hostname, conf.Port, auth}
}

func (c *Client) Send(msg *Message) (io.Reader, error) {
	logger.Debug().Msgf("sending %#v", msg)

	var h mail.Header
	h.SetMessageID(c.generateMessageId())
	h.Set("User-agent", userAgent)
	h.SetDate(time.Now())
	h.SetSubject(msg.Subject)
	h.SetAddressList("From", []*mail.Address{c.address})

	to := make([]*mail.Address, 0)
	for _, addr := range msg.To {
		parsed, err := mail.ParseAddress(addr)
		if err == nil {
			to = append(to, parsed)
		}
	}
	h.SetAddressList("To", to)

	if len(msg.CC) > 0 {
		cc := make([]*mail.Address, 0)
		for _, addr := range msg.CC {
			parsed, err := mail.ParseAddress(addr)
			if err == nil {
				cc = append(cc, parsed)
			}
		}
		h.SetAddressList("Cc", cc)
	}

	var buf bytes.Buffer
	mw, err := mail.CreateWriter(&buf, h)
	if err != nil {
		return nil, fmt.Errorf("mail: %w", err)
	}

	var th mail.InlineHeader
	th.Set("Content-Type", "text/plain")
	tw, err := mw.CreateInline()
	if err != nil {
		return nil, fmt.Errorf("mail body: %w", err)
	}
	w, err := tw.CreatePart(th)
	if err != nil {
		return nil, fmt.Errorf("mail body: %w", err)
	}
	if _, err := io.WriteString(w, msg.Body); err != nil {
		return nil, fmt.Errorf("mail body: %w", err)
	}
	w.Close()
	tw.Close()

	for _, attach := range msg.Attachments {
		var ah mail.AttachmentHeader
		ah.SetFilename(attach.Filename)

		w, err = mw.CreateAttachment(ah)
		if err != nil {
			return nil, fmt.Errorf("mail attachment: %w", err)
		}
		if _, err := io.Copy(w, attach); err != nil {
			return nil, fmt.Errorf("mail attachment: %w", err)
		}
		w.Close()
	}

	pass, err := c.auth.Password(c.username, c.hostname)
	if err != nil {
		return nil, fmt.Errorf("auth: %w", err)
	}

	auth := sasl.NewPlainClient("", c.username, pass)
	toAddr := make([]string, len(to))
	for idx, addr := range to {
		toAddr[idx] = addr.Address
	}

	var out bytes.Buffer
	return &out, smtp.SendMail(fmt.Sprintf("%s:%d", c.hostname, c.port), auth, c.address.Address, toAddr, io.TeeReader(&buf, &out))
}

func (c *Client) generateMessageId() string {
	items := strings.Split(c.address.Address, "@")

	id := make([]byte, 5)
	rand.Read(id)

	return fmt.Sprintf("%s.amail@%s", hex.EncodeToString(id), items[1])
}
