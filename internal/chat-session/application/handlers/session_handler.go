package handlers

import (
	"context"
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/application/service"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/model"
	"strings"
)

type CreateSessionRequest struct {
	SessionID string `json:"session_id"`
	Name      string `json:"name"`
}

func (r *CreateSessionRequest) Validate() error {
	r.SessionID = strings.TrimSpace(r.SessionID)
	r.Name = strings.TrimSpace(r.Name)

	if r.SessionID == "" {
		return errors.New("session_id is required")
	} else if r.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

type CreateSessionHandler struct {
	Command  CreateSessionRequest
	Response *model.Session
	*service.SessionService
}

func (h *CreateSessionHandler) Handle(ctx context.Context) error {
	session, err := h.SessionService.CreateSession(ctx, h.Command.SessionID, h.Command.Name)
	if err != nil {
		return err
	}
	h.Response = session
	return nil
}
