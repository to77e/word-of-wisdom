package client

import (
	_ "encoding/json"
)

// MessageType describes the type of message
// and is used to validate the message type in the request
// and is necessary for distinguish a regular message from an error message
type MessageType int

const (
	ChallengeRequest MessageType = iota + 1
	ChallengeResponse
	SolutionRequest
	SolutionResponse
	Error
)

type ChallengeMessageRequest struct {
	Type    MessageType `json:"type" validate:"required,eq=1"`
	Content string      `json:"content" validate:"required"`
}

type ChallengeMessageResponse struct {
	Type         MessageType `json:"type" validate:"required,eq=2|eq=5"`
	Challenge    []byte      `json:"challenge" validate:"required_if=Type 2"`
	ErrorMessage string      `json:"error_message" validate:"required_if=Type 5"`
}

type SolutionMessageRequest struct {
	Type     MessageType `json:"type" validate:"required,eq=3"`
	Solution []byte      `json:"solution" validate:"required"`
}

type SolutionMessageResponse struct {
	Type         MessageType `json:"type" validate:"required,eq=4|eq=5"`
	Result       string      `json:"result" validate:"required_if=Type 4"`
	ErrorMessage string      `json:"error_message" validate:"required_if=Type 5"`
}
