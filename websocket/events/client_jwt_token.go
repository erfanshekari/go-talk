package events

type ClientJWTToken struct {
	AccessToken string `json:"accessToken" bson:"accessToken"`
}
