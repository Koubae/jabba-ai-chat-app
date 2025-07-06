package handlers

import (
	"context"
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/application/service"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/model"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/settings"
	"strings"
)

type CreateSessionRequest struct {
	SessionID string `json:"session_id"`
	Name      string `json:"name"`
	MemberID  string `json:"member_id" binding:"required"`
	Channel   string `json:"channel" binding:"required"`
}

func (r *CreateSessionRequest) Validate() error {
	r.SessionID = strings.TrimSpace(r.SessionID)
	r.Name = strings.TrimSpace(r.Name)
	r.MemberID = strings.TrimSpace(r.MemberID)
	r.Channel = strings.TrimSpace(r.Channel)

	if r.SessionID == "" {
		return errors.New("session_id is required")
	} else if r.Name == "" {
		return errors.New("name is required")
	} else if r.MemberID == "" {
		return errors.New("member_id is required")
	} else if r.Channel == "" {
		return errors.New("channel is required")
	}

	return nil
}

type CreateSessionResponse struct {
	ChatURL string `json:"chat_url"`
	model.Session
}

type CreateSessionHandler struct {
	Command  CreateSessionRequest
	Response *CreateSessionResponse
	*service.SessionService
}

func (h *CreateSessionHandler) Handle(ctx context.Context) error {
	session, err := h.SessionService.CreateSession(
		ctx,
		h.Command.SessionID,
		h.Command.Name,
		h.Command.MemberID,
		h.Command.Channel,
	)
	if err != nil {
		return err
	}

	config := settings.GetConfig()
	h.Response = &CreateSessionResponse{
		ChatURL: config.GetURL(),
		Session: *session,
	}
	return nil
}
