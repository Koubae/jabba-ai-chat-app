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

type Session struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	ApplicationID string    `bson:"application_id,omitempty" json:"application_id"`
	UserID        string    `bson:"user_id,omitempty" json:"user_id"`
	Name          string    `json:"name"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
}
