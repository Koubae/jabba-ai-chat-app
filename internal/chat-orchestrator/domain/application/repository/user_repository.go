package repository

import (
	"context"
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, applicationID string, IdentityID int64) (*model.User, error)
}

var (
	ErrUserAlreadyExists = errors.New("USER_ALREADY_EXISTS")
	ErrUserOnCreate      = errors.New("USER_ERROR_ON_CREATE")
	ErrUserNotFound      = errors.New("USER_NOT_FOUND")
)
