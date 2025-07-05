package handlers

import (
	"context"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/application/service"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
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
