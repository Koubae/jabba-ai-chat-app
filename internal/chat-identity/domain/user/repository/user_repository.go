package repository

import (
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/domain/user/model"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id int64) (*model.User, error)
	GetByUsername(applicationID string, username string) (*model.User, error)
}

var (
	ErrUserNotFound         = errors.New("USER_NOT_FOUND")
	ErrUserAlreadyExists    = errors.New("USER_ALREADY_EXISTS")
	ErrUserIdentityMismatch = errors.New("USER_IDENTITY_MISMATCH")
)
