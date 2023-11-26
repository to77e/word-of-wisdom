package server

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
	Type      MessageType `json:"type" validate:"required,eq=2"`
	Challenge []byte      `json:"challenge" validate:"required"`
}

type SolutionMessageRequest struct {
	Type     MessageType `json:"type" validate:"required,eq=3"`
	Solution []byte      `json:"solution" validate:"required"`
}

type SolutionMessageResponse struct {
	Type   MessageType `json:"type" validate:"required,eq=4"`
	Result string      `json:"result" validate:"required"`
}

type ErrorMessage struct {
	Type         MessageType `json:"type" validate:"required,eq=5"`
	ErrorMessage string      `json:"error_message" validate:"required"`
}
