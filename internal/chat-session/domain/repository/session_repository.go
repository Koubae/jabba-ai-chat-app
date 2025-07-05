package repository

import (
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/model"
)

type SessionRepository interface {
	Create(session *model.Session) error
	Get(applicationID string, sessionID string) (*model.Session, error)
}

var (
	ErrSessionNotFound = errors.New("SESSION_NOT_FOUND")
	ErrSessionParse    = errors.New("SESSION_ERROR_PARSING")
)
