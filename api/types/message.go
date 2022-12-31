package types

type MessageType string

const (
	Text MessageType = "text"
	File MessageType = "file"
)

type Message struct {
	Type MessageType `json:"type" validate:"required"`
}
