package handlers

import (
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/application/service"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
	"golang.org/x/net/context"
)

type CreateSessionRequest struct {
	ApplicationID string `json:"application_id"`
	IdentityID    int64  `json:"identity_id"`
	Name          string `json:"name"`
}

type CreateSessionHandler struct {
	Command  CreateSessionRequest
	Response *model.Session
	*service.SessionService
}

func (h *CreateSessionHandler) Handle(ctx context.Context) error {
	session, err := h.SessionService.Create(ctx, h.Command.ApplicationID, h.Command.IdentityID, h.Command.Name)
	if err != nil {
		return err
	}
	h.Response = session
	return nil
}
