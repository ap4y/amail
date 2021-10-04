package notmuch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	notmuchBin    = "notmuch"
	formatVersion = "4"
)

type CountOutputType string
type ReplyToType string

const (
	CountOutputMessages CountOutputType = "messages"
	CountOutputThreads  CountOutputType = "threads"
	CountOutputFiles    CountOutputType = "files"

	ReplyToAll    ReplyToType = "all"
	ReplyToSender ReplyToType = "sender"
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

func (c *Client) SearchFiles(term string) ([]string, error) {
	var files []string
	if err := c.jsonExec(&files, "search", "--output=files", term); err != nil {
		return nil, err
	}

	return files, nil
}

func (c *Client) Show(term string) ([][]interface{}, error) {
	var messages [][]interface{}
	if err := c.jsonExec(&messages, "show", term); err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *Client) Attachment(messageID, partID string) (io.ReadSeeker, map[string]interface{}, error) {
	args := []string{"--part", partID, "id:" + messageID}

	var part map[string]interface{}
	if err := c.jsonExec(&part, "show", args...); err != nil {
		return nil, nil, err
	}

	res, err := c.exec("show", nil, args...)
	if err != nil {
		return nil, nil, err
	}

	return bytes.NewReader(res), part, nil
}

func (c *Client) Count(term string, output CountOutputType) (int, error) {
	out, err := c.exec("count", nil, "--output", string(output), term)
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
	_, err := c.exec("tag", nil, args...)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Dump(term string) ([]string, error) {
	res, err := c.exec("dump", nil, term)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`[+](\w*)\s`)
	matches := re.FindAllSubmatch(res, -1)
	tags := make([]string, len(matches))
	for idx, match := range re.FindAllSubmatch(res, -1) {
		tags[idx] = string(match[1])
	}
	return tags, nil
}

func (c *Client) Reply(term string, replyTo ReplyToType) (*Reply, error) {
	args := []string{"--reply-to", string(replyTo), term}

	var reply *Reply
	if err := c.jsonExec(&reply, "reply", args...); err != nil {
		return nil, err
	}

	return reply, nil
}

func (c *Client) Index() error {
	_, err := c.exec("new", nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Insert(folder string, msg io.Reader, tags ...string) error {
	args := []string{"--folder", folder, "--keep"}
	args = append(args, tags...)

	_, err := c.exec("insert", msg, args...)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Address(term string) ([]map[string]string, error) {
	var addresses []map[string]string
	if err := c.jsonExec(&addresses, "address", term); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (c *Client) exec(cmd string, in io.Reader, args ...string) ([]byte, error) {
	allArgs := append([]string{cmd}, args...)
	ec := exec.Command(c.binPath, allArgs...)
	if in != nil {
		ec.Stdin = in
	}

	out, err := ec.Output()
	if err != nil {
		return nil, fmt.Errorf("exec: %w. %s", err, out)
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
