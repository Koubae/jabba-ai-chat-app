package model

import "time"

type Application struct {
	ID      string    `json:"id" db:"id"`
	Name    string    `json:"name" db:"name"`
	Created time.Time `json:"created" db:"created"`
	Updated time.Time `json:"updated" db:"updated"`
}
