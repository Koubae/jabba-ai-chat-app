package repository

import (
	"context"
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/model"
)

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session, identityID int64) error
	Get(ctx context.Context, applicationID string, sessionID string, identityID int64) (*model.Session, error)
}

var (
	ErrSessionNotFound = errors.New("SESSION_NOT_FOUND")
	ErrSessionParse    = errors.New("SESSION_ERROR_PARSING")
)
