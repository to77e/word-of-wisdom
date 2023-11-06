package message

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Message struct {
	data      []byte
	validator *validator.Validate
}

func New(data []byte, validate *validator.Validate) *Message {
	return &Message{
		data:      data,
		validator: validate,
	}
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) UnmarshalData(data interface{}) error {
	if err := json.Unmarshal(m.data, data); err != nil {
		return fmt.Errorf("unmarshal data: %w\n", err)
	}
	if err := m.validator.Struct(data); err != nil {
		return fmt.Errorf("validate data: %w\n", err)
	}
	return nil
}

func (m *Message) MarshalData(data interface{}) error {
	var err error
	m.data, err = json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data: %w\n", err)
	}
	return nil
}
