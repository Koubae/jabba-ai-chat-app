package model

type User struct {
	ID            uint   `json:"id"`
	ApplicationID string `json:"application_id"`
	Username      string `json:"username"`
}
