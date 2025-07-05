package repository

import (
	"context"
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
)

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error
	GetSession(ctx context.Context, applicationID string, userID string, name string) (*model.Session, error)
	GetSessionByID(ctx context.Context, sessionID string) (*model.Session, error)
	ListWithPagination(ctx context.Context, applicationID string, userID string, limit int64, offset int64) ([]*model.Session, error)
}

var (
	ErrSessionAlreadyExists = errors.New("SESSION_ALREADY_EXISTS")
	ErrSessionOnCreate      = errors.New("SESSION_ERROR_ON_CREATE")
	ErrSessionNotFound      = errors.New("SESSION_NOT_FOUND")
)
