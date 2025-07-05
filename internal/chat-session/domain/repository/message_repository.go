package repository

import (
	"context"
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/model"
)

type MessageRepository interface {
	AddMessage(ctx context.Context, applicationID string, sessionID string, message *model.Message) error
	GetMessages(ctx context.Context, applicationID string, sessionID string) ([]*model.Message, error)
}
