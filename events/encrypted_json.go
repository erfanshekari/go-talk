package events

type EncryptedJsonType string

const (
	Byte      EncryptedJsonType = "byte"
	ByteArray EncryptedJsonType = "byteArray"
)

type EncryptedJson struct {
	Type    EncryptedJsonType `json:"type" bson:"type"`
	Content any               `json:"content" bson:"content"`
}
