package client

import (
	"fmt"
	"net"

	"github.com/to77e/word-of-wisdom/internal/message"
	"github.com/to77e/word-of-wisdom/tools/validator"

	_ "github.com/to77e/word-of-wisdom/internal/message"
)

const (
	defaultBufferSize = 1024
)

type Client struct {
	connection net.Conn
}

func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("dial tcp connection: %w\n", err)
	}
	return &Client{connection: conn}, nil
}

func (c *Client) GetRemoteAddress() string {
	return c.connection.RemoteAddr().String()
}

func (c *Client) Close() error {
	if err := c.connection.Close(); err != nil {
		return fmt.Errorf("close tcp connection: %w\n", err)
	}
	return nil
}

func (c *Client) WriteChallengeRequest() error {
	challengeRequest := &ChallengeMessageRequest{
		Type:    ChallengeRequest,
		Content: "start",
	}
	if err := c.sendRequest(challengeRequest); err != nil {
		return fmt.Errorf("send challenge request: %w\n", err)
	}

	return nil
}

func (c *Client) ReadChallengeResponse() ([]byte, error) {
	data := &ChallengeMessageResponse{
		Challenge: make([]byte, 0),
	}
	if err := c.getResponse(data); err != nil {
		return nil, fmt.Errorf("get challenge response: %w\n", err)
	}

	if data.Type == Error {
		return nil, fmt.Errorf("error message: %s\n", data.ErrorMessage)
	}

	return data.Challenge, nil
}

func (c *Client) WriteSolutionRequest(solution []byte) error {
	solutionRequest := &SolutionMessageRequest{
		Type:     SolutionRequest,
		Solution: solution,
	}
	if err := c.sendRequest(solutionRequest); err != nil {
		return fmt.Errorf("send solution request: %w\n", err)
	}

	return nil
}

func (c *Client) ReadSolutionResponse() (string, error) {
	data := &SolutionMessageResponse{}
	if err := c.getResponse(data); err != nil {
		return "", fmt.Errorf("get solution response: %w\n", err)
	}

	if data.Type == Error {
		return "", fmt.Errorf("error message: %s\n", data.ErrorMessage)
	}

	return data.Result, nil
}

func (c *Client) getResponse(resp interface{}) error {
	// todo: dynamic buffer size
	buf := make([]byte, defaultBufferSize)
	n, err := c.connection.Read(buf)
	if err != nil {
		return fmt.Errorf("read tcp connection: %w\n", err)
	}

	m := message.New(buf[:n:n], validator.Get())
	if err = m.UnmarshalData(resp); err != nil {
		return fmt.Errorf("unmarshal data: %w\n", err)
	}

	return nil
}

func (c *Client) sendRequest(data interface{}) error {
	m := message.New(nil, nil)

	if err := m.MarshalData(data); err != nil {
		return fmt.Errorf("marshal data: %w\n", err)
	}

	if _, err := c.connection.Write(m.GetData()); err != nil {
		return fmt.Errorf("write tcp connection: %w\n", err)
	}

	return nil
}
