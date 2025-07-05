package repository

import (
	"context"
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
)

type ApplicationRepository interface {
	Create(ctx context.Context, application *model.Application) error
	GetByID(id int64) (*model.Application, error)
	GetByName(name string) (*model.Application, error)
}

var (
	ErrApplicationOnCreate         = errors.New("APPLICATION_ERROR_ON_CREATE")
	ErrApplicationNotFound         = errors.New("APPLICATION_NOT_FOUND")
	ErrApplicationAlreadyExists    = errors.New("APPLICATION_ALREADY_EXISTS")
	ErrApplicationIdentityMismatch = errors.New("APPLICATION_IDENTITY_MISMATCH")
)
