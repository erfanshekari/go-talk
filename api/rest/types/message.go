package types

type MessageType string

const (
	Text MessageType = "text"
	File MessageType = "file"
)

type Message struct {
	Type    MessageType `json:"type" validate:"required"`
	Channel string      `json:"channel" validate:"required"`
	From    string      `json:"from" validate:"required"`
	Content string      `json:"content" validate:"required"`
}
