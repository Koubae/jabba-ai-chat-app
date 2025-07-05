package model

import "time"

type Session struct {
	ID            string     `json:"id"`
	ApplicationID string     `json:"application_id"`
	Name          string     `json:"name"`
	Created       *time.Time `json:"created"`
	Updated       *time.Time `json:"updated"`
}
