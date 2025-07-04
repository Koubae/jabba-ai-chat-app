package model

import "time"

type User struct {
	ID            int64     `json:"id" db:"id"`
	ApplicationID string    `json:"application_id" db:"application_id"`
	Username      string    `json:"username" db:"username"`
	PasswordHash  string    `json:"password_hash" db:"password_hash"`
	Created       time.Time `json:"created" db:"created"`
	Updated       time.Time `json:"updated" db:"updated"`
}
