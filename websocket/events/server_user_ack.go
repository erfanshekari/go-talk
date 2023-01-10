package events

type ServerUserACK struct {
	UserID string `json:"user_id" bson:"user_id"`
}
