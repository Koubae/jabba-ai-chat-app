package handlers

import (
	"context"
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/application/service"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

type CreateApplicationRequest struct {
	Name string `json:"name"`
}

type CreateApplicationHandler struct {
	Command  CreateApplicationRequest
	Response *model.Application
	*service.ApplicationService
}

func (h *CreateApplicationHandler) Handle(ctx context.Context) error {
	application, err := h.ApplicationService.Create(ctx, h.Command.Name)
	if err != nil {
		return err
	}
	h.Response = application
	return nil
}

type GetApplicationRequest struct {
	Name string `json:"name"`
}

type GetApplicationHandler struct {
	Command  GetApplicationRequest
	Response *model.Application
	*service.ApplicationService
}

func (h *GetApplicationHandler) Handle(ctx context.Context) error {
	application, err := h.ApplicationService.Get(ctx, h.Command.Name)
	if err != nil {
		return err
	}
	h.Response = application
	return nil
}

type ListApplicationRequest struct {
	Limit  int64
	Offset int64
}

func (r *ListApplicationRequest) Validate(c *gin.Context) error {
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

type ListApplicationHandler struct {
	Command  ListApplicationRequest
	Response []*model.Application
	*service.ApplicationService
}

func (h *ListApplicationHandler) Handle(ctx context.Context) error {
	applications, err := h.ApplicationService.List(ctx, h.Command.Limit, h.Command.Offset)
	if err != nil {
		return err
	}
	h.Response = applications
	return nil
}
