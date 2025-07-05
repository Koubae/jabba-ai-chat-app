package model

import "time"

type Session struct {
	ID            string    `json:"id"`
	ApplicationID string    `json:"application_id"`
	Name          string    `json:"name"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
}

type Message struct {
	ApplicationID string `json:"application_id"`
	SessionID     string `json:"session_id"`
	Role          string `json:"role"`
	UserID        int    `json:"user_id"`
	Username      string `json:"username"`
	Message       string `json:"message"`
	Timestamp     int64  `json:"timestamp"`
	Member
}

type Member struct {
	MemberID string `json:"member_id"`
	Channel  string `json:"channel"`
}
