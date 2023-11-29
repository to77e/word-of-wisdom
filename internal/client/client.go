package client

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/to77e/word-of-wisdom/internal/messages"
	"github.com/to77e/word-of-wisdom/tools/validator"
)

type Client struct {
	connection net.Conn
}

type Challenge struct {
	Content    []byte
	Difficulty int
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
	challengeRequest := &messages.ChallengeMessageRequest{
		Type:    messages.ChallengeRequest,
		Content: "start",
	}
	if err := c.sendRequest(challengeRequest); err != nil {
		return fmt.Errorf("send challenge request: %w\n", err)
	}

	return nil
}

func (c *Client) ReadChallengeResponse() (*Challenge, error) {
	data := &messages.ChallengeMessageResponse{
		Challenge: make([]byte, 0),
	}
	if err := c.getResponse(data); err != nil {
		return nil, fmt.Errorf("get challenge response: %w\n", err)
	}

	if data.Type == messages.Error {
		return nil, fmt.Errorf("error messages: %s\n", data.ErrorMessage)
	}

	return &Challenge{
		Content:    data.Challenge,
		Difficulty: data.Difficulty,
	}, nil
}

func (c *Client) WriteSolutionRequest(solution []byte) error {
	solutionRequest := &messages.SolutionMessageRequest{
		Type:     messages.SolutionRequest,
		Solution: solution,
	}
	if err := c.sendRequest(solutionRequest); err != nil {
		return fmt.Errorf("send solution request: %w\n", err)
	}

	return nil
}

func (c *Client) ReadSolutionResponse() (string, error) {
	data := &messages.SolutionMessageResponse{}
	if err := c.getResponse(data); err != nil {
		return "", fmt.Errorf("get solution response: %w\n", err)
	}

	if data.Type == messages.Error {
		return "", fmt.Errorf("error messages: %s\n", data.ErrorMessage)
	}

	return data.Quote, nil
}

func (c *Client) getResponse(resp interface{}) error {
	decoder := json.NewDecoder(c.connection)
	if err := decoder.Decode(resp); err != nil {
		return fmt.Errorf("decode response: %w\n", err)
	}

	if err := validator.GetInstance().Struct(resp); err != nil {
		return fmt.Errorf("validate data: %w\n", err)
	}

	return nil
}

func (c *Client) sendRequest(data interface{}) error {
	encoder := json.NewEncoder(c.connection)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encode data: %w\n", err)
	}

	return nil
}
