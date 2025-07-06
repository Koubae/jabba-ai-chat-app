package model

import (
	"time"
)

type Application struct {
	ID      string    `bson:"_id,omitempty" json:"id"`
	Name    string    `json:"name" db:"name"`
	Created time.Time `json:"created" db:"created"`
	Updated time.Time `json:"updated" db:"updated"`
}

type User struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	ApplicationID string    `bson:"application_id,omitempty" json:"application_id"`
	IdentityID    int64     `bson:"identity_id"`
	Username      string    `json:"username"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
}

type Session struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	ApplicationID string    `bson:"application_id,omitempty" json:"application_id"`
	UserID        string    `bson:"user_id,omitempty" json:"user_id"`
	Name          string    `json:"name"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
}

type SessionConnection struct {
	ChatURL string `json:"chat_url"`
	Session
	UserID string  `json:"user_id,omitempty"`
	Owner  *Member `json:"owner"`
}

type Member struct {
	Role     string `json:"role"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	MemberID string `json:"member_id"`
	Channel  string `json:"channel"`
}
