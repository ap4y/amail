package notmuch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

const (
	notmuchBin    = "notmuch"
	formatVersion = "4"
)

type Client struct {
	binPath string
}

func NewClient() (*Client, error) {
	path, err := exec.LookPath(notmuchBin)
	if err != nil {
		return nil, fmt.Errorf("failed to find notmuch binary: %s", err)
	}

	return &Client{path}, nil
}

func (c *Client) Search(term string, limit, offset int) ([]Thread, error) {
	args := []string{}
	if limit > 0 {
		args = append(args, "--limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		args = append(args, "--offset", strconv.Itoa(offset))
	}
	args = append(args, term)

	var threads []Thread
	if err := c.jsonExec(&threads, "search", args...); err != nil {
		return nil, err
	}

	return threads, nil
}

func (c *Client) Show(term string) ([][]interface{}, error) {
	var messages [][]interface{}
	if err := c.jsonExec(&messages, "show", "--include-html", term); err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *Client) Attachment(messageID, partID string) (io.ReadSeeker, error) {
	res, err := c.exec("show", "--part", partID, "id:"+messageID)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(res), nil
}

func (c *Client) Count(term string) (int, error) {
	out, err := c.exec("count", term)
	if err != nil {
		return -1, err
	}

	cnt, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return -1, err
	}

	return cnt, nil
}

func (c *Client) Tag(term string, tags []string) error {
	args := append(tags, "--", term)
	_, err := c.exec("tag", args...)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Reply(term, replyTo string) (*Reply, error) {
	args := []string{"--reply-to", replyTo, term}

	var reply *Reply
	if err := c.jsonExec(&reply, "reply", args...); err != nil {
		return nil, err
	}

	return reply, nil
}

func (c *Client) Index() error {
	_, err := c.exec("new")
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) exec(cmd string, args ...string) ([]byte, error) {
	allArgs := append([]string{cmd}, args...)
	ec := exec.Command(c.binPath, allArgs...)

	out, err := ec.Output()
	if err != nil {
		return nil, fmt.Errorf("exec: %v", err)
	}

	return out, nil
}

func (c *Client) jsonExec(res interface{}, cmd string, args ...string) error {
	allArgs := append([]string{cmd, "--format-version", formatVersion, "--format", "json"}, args...)
	ec := exec.Command(c.binPath, allArgs...)
	stdout, err := ec.StdoutPipe()
	if err != nil {
		return fmt.Errorf("exec: %v", err)
	}
	stderr, err := ec.StderrPipe()
	if err != nil {
		return fmt.Errorf("exec: %v", err)
	}

	if err := ec.Start(); err != nil {
		return fmt.Errorf("exec: %v.", err)
	}

	if err := json.NewDecoder(stdout).Decode(res); err != nil {
		out, _ := io.ReadAll(stderr)
		return fmt.Errorf("json: %v. %s", err, out)
	}

	if err := ec.Wait(); err != nil {
		return fmt.Errorf("exec: %v", err)
	}

	return nil
}
