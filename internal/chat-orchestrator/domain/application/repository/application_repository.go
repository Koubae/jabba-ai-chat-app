package repository

import (
	"context"
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
)

type ApplicationRepository interface {
	Create(ctx context.Context, application *model.Application) error
	GetByID(ctx context.Context, id string) (*model.Application, error)
	GetByName(ctx context.Context, name string) (*model.Application, error)
}

var (
	ErrApplicationAlreadyExists = errors.New("APPLICATION_ALREADY_EXISTS")
	ErrApplicationOnCreate      = errors.New("APPLICATION_ERROR_ON_CREATE")
	ErrApplicationNotFound      = errors.New("APPLICATION_NOT_FOUND")
)
