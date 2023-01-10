package events

type ClientPublicKey struct {
	PublicKey string `json:"publicKey"`
}

type ClientJWTToken struct {
	AccessToken string `json:"accessToken" bson:"accessToken"`
}

type ServerPublicKey struct {
	PublicKey string `json:"publicKey"`
}

type ServerUserACK struct {
	UserID string `json:"userID" bson:"userID"`
}
