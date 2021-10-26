package smtp

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"ap4y.me/cloud-mail/config"
	"github.com/emersion/go-message/mail"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog"
)

var ErrInvalidDkimKey = errors.New("invalid DKIM key")

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
	DkimKey(username, hostname string) (string, error)
}

type Headers map[string]string

type Attachment struct {
	io.ReadCloser
	Filename    string
	ContentType string
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

func (c *Client) Compose(msg *Message) (io.Reader, []*mail.Address, []*mail.Address, error) {
	var h mail.Header
	h.SetMessageID(c.generateMessageId())
	h.Set("User-agent", userAgent)
	h.SetDate(time.Now())
	h.SetSubject(msg.Subject)
	h.SetAddressList("From", []*mail.Address{c.address})
	for key, value := range msg.Headers {
		h.Set(key, value)
	}

	to := make([]*mail.Address, 0)
	for _, addr := range msg.To {
		parsed, err := mail.ParseAddress(addr)
		if err == nil {
			to = append(to, parsed)
		}
	}
	h.SetAddressList("To", to)

	var cc []*mail.Address
	if len(msg.CC) > 0 {
		cc = make([]*mail.Address, 0)
		for _, addr := range msg.CC {
			parsed, err := mail.ParseAddress(addr)
			if err == nil {
				cc = append(cc, parsed)
			}
		}
		h.SetAddressList("Cc", cc)
	}

	var buf bytes.Buffer
	if len(msg.Attachments) == 0 {
		h.Set("Content-Type", "text/plain")
		w, err := mail.CreateSingleInlineWriter(&buf, h)
		if err != nil {
			return nil, to, cc, fmt.Errorf("mail: %w", err)
		}

		if _, err := io.WriteString(w, msg.Body); err != nil {
			return nil, to, cc, fmt.Errorf("mail body: %w", err)
		}
		w.Close()
	} else {
		mw, err := mail.CreateWriter(&buf, h)
		if err != nil {
			return nil, to, cc, fmt.Errorf("mail: %w", err)
		}

		var th mail.InlineHeader
		th.Set("Content-Type", "text/plain")
		tw, err := mw.CreateInline()
		if err != nil {
			return nil, to, cc, fmt.Errorf("mail body: %w", err)
		}
		w, err := tw.CreatePart(th)
		if err != nil {
			return nil, to, cc, fmt.Errorf("mail body: %w", err)
		}
		if _, err := io.WriteString(w, msg.Body); err != nil {
			return nil, to, cc, fmt.Errorf("mail body: %w", err)
		}
		w.Close()
		tw.Close()

		for _, attach := range msg.Attachments {
			var ah mail.AttachmentHeader
			ah.Set("Content-Type", attach.ContentType)
			ah.SetFilename(attach.Filename)

			w, err = mw.CreateAttachment(ah)
			if err != nil {
				return nil, to, cc, fmt.Errorf("mail attachment: %w", err)
			}
			if _, err := io.Copy(w, attach); err != nil {
				return nil, to, cc, fmt.Errorf("mail attachment: %w", err)
			}
			w.Close()
		}

		mw.Close()
	}

	return &buf, to, cc, nil
}
func (c *Client) Send(msg *Message) (io.Reader, error) {
	logger.Debug().Msgf("sending %#v", msg)

	m, to, cc, err := c.Compose(msg)
	if err != nil {
		return nil, err
	}

	if signed, err := c.dkimSign(m); err != nil {
		return nil, fmt.Errorf("dkim: %w", err)
	} else {
		m = signed
	}

	pass, err := c.auth.Password(c.username, c.hostname)
	if err != nil {
		return nil, fmt.Errorf("auth: %w", err)
	}

	auth := sasl.NewPlainClient("", c.username, pass)
	toAddr := make([]string, len(to)+len(cc))
	for idx, addr := range to {
		toAddr[idx] = addr.Address
	}
	for idx, addr := range cc {
		toAddr[len(to)+idx] = addr.Address
	}

	var out bytes.Buffer
	return &out, smtp.SendMail(
		fmt.Sprintf("%s:%d", c.hostname, c.port),
		auth,
		c.address.Address,
		toAddr,
		io.TeeReader(m, &out),
	)
}

func (c *Client) getDomain() string {
	items := strings.Split(c.address.Address, "@")
	return items[1]
}

func (c *Client) generateMessageId() string {
	id := make([]byte, 5)
	rand.Read(id)

	return fmt.Sprintf("%s.amail@%s", hex.EncodeToString(id), c.getDomain())
}

func (c *Client) dkimSign(msg io.Reader) (io.Reader, error) {
	dkimKey, err := c.auth.DkimKey(c.username, c.hostname)
	if err != nil {
		return msg, nil
	}

	privKey, err := getPrivKey([]byte(dkimKey))
	if err != nil {
		return nil, err
	}

	opts := &dkim.SignOptions{
		Domain:                 c.getDomain(),
		Selector:               "default",
		Signer:                 privKey,
		HeaderCanonicalization: dkim.CanonicalizationRelaxed,
		BodyCanonicalization:   dkim.CanonicalizationRelaxed,
	}

	var buf bytes.Buffer
	return &buf, dkim.Sign(&buf, msg, opts)
}

func getPrivKey(pemData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, ErrInvalidDkimKey
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
