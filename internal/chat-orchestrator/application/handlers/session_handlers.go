package handlers

import (
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/application/service"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"strconv"
)

type CreateSessionRequest struct {
	ApplicationID string `json:"application_id"`
	IdentityID    int64  `json:"identity_id"`
	Username      string `json:"username"`
	Name          string `json:"name"`
}

type CreateSessionHandler struct {
	Command  CreateSessionRequest
	Response *model.Session
	*service.SessionService
}

func (h *CreateSessionHandler) Handle(ctx context.Context) error {
	session, err := h.SessionService.Create(ctx, h.Command.ApplicationID, h.Command.IdentityID, h.Command.Username, h.Command.Name)
	if err != nil {
		return err
	}
	h.Response = session
	return nil
}

type GetSessionRequest struct {
	ApplicationID string `json:"application_id"`
	IdentityID    int64  `json:"identity_id"`
	Name          string `json:"name"`
}

type GetSessionHandler struct {
	Command  GetSessionRequest
	Response *model.Session
	*service.SessionService
}

func (h *GetSessionHandler) Handle(ctx context.Context) error {
	session, err := h.SessionService.Get(ctx, h.Command.ApplicationID, h.Command.IdentityID, h.Command.Name)
	if err != nil {
		return err
	}
	h.Response = session
	return nil
}

type ListSessionRequest struct {
	ApplicationID string `json:"application_id"`
	IdentityID    int64  `json:"identity_id"`
	Limit         int64
	Offset        int64
}

func (r *ListSessionRequest) Validate(c *gin.Context) error {
	// Default values
	limit := int64(100)
	offset := int64(0)

	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, parseErr := strconv.ParseInt(limitStr, 10, 64)
		if parseErr != nil {
			return errors.New("invalid limit parameter: must be a number")
		}
		limit = parsedLimit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		parsedOffset, parseErr := strconv.ParseInt(offsetStr, 10, 64)
		if parseErr != nil {
			return errors.New("invalid offset parameter: must be a number")
		}
		offset = parsedOffset
	}

	if limit < 0 {
		return errors.New("limit cannot be negative")
	}
	if limit > 10000 {
		return errors.New("limit cannot exceed 10000")
	}

	if offset < 0 {
		return errors.New("offset cannot be negative")
	}

	r.Limit = limit
	r.Offset = offset
	return nil
}

type ListSessionHandler struct {
	Command  ListSessionRequest
	Response []*model.Session
	*service.SessionService
}

func (h *ListSessionHandler) Handle(ctx context.Context) error {
	entities, err := h.SessionService.List(
		ctx,
		h.Command.ApplicationID,
		h.Command.IdentityID,
		h.Command.Limit,
		h.Command.Offset,
	)
	if err != nil {
		return err
	}
	h.Response = entities
	return nil
}
