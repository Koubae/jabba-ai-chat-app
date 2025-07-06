package model

import (
	"errors"
	"time"
)

type Message struct {
	ApplicationID string `json:"application_id"`
	SessionID     string `json:"session_id"`
	Message       string `json:"message"`
	Timestamp     int64  `json:"timestamp"`
	Member
}

type Member struct {
	Role     string `json:"role"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	MemberID string `json:"member_id"`
	Channel  string `json:"channel"`
}

type Session struct {
	ID            string    `json:"id"`
	ApplicationID string    `json:"application_id"`
	Name          string    `json:"name"`
	Owner         *Member   `json:"owner"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
}

func (s *Session) IsSameOwner(other *Session) bool {
	if s.Owner == nil && other.Owner == nil {
		return true
	}
	if s.Owner == nil || other.Owner == nil {
		return false
	}
	return s.Owner.UserID == other.Owner.UserID
}

var ErrIsNotOwnerOfSession = errors.New("MEMBER_IS_NOT_OWNER_OF_SESSION")
