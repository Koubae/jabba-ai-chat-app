package model

type User struct {
	UserID        uint   `json:"user_id"`
	Username      string `json:"username"`
	ApplicationID string `json:"application_id"`
}
