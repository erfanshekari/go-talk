package events

type BytesWrappedJsonType string

const (
	Byte      BytesWrappedJsonType = "byte"
	ByteArray BytesWrappedJsonType = "byteArray"
)

type BytesWrappedJson struct {
	Type    BytesWrappedJsonType `json:"type" bson:"type"`
	Content any                  `json:"content" bson:"content"`
}
