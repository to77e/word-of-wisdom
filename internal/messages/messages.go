package messages

// Type describes the type of messages
// and is used to validate the messages type in the request
// and is necessary for distinguish a regular messages from an error messages
type Type int

const (
	ChallengeRequest Type = iota + 1
	ChallengeResponse
	SolutionRequest
	SolutionResponse
	Error
)

type ChallengeMessageRequest struct {
	Type    Type   `json:"type" validate:"required,eq=1"`
	Content string `json:"content"`
}

type ChallengeMessageResponse struct {
	Type         Type   `json:"type" validate:"required,eq=2|eq=5"`
	Challenge    []byte `json:"challenge" validate:"required_if=Type 2"`
	Difficulty   int    `json:"difficulty" validate:"required_if=Type 2"`
	ErrorMessage string `json:"error_message" validate:"required_if=Type 5"`
}

type SolutionMessageRequest struct {
	Type     Type   `json:"type" validate:"required,eq=3"`
	Solution []byte `json:"solution" validate:"required"`
}

type SolutionMessageResponse struct {
	Type         Type   `json:"type" validate:"required,eq=4|eq=5"`
	Quote        string `json:"quote" validate:"required_if=Type 4"`
	ErrorMessage string `json:"error_message" validate:"required_if=Type 5"`
}

type ErrorMessage struct {
	Type         Type   `json:"type" validate:"required,eq=5"`
	ErrorMessage string `json:"error_message" validate:"required"`
}
