package handlers

import (
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/application/service"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/domain/user/model"
	domainrepository "github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/domain/user/repository"
)

type GetAccountRequest struct {
	ApplicationID string
	Username      string
}

type GetAccountResponse struct {
	*model.UserResponse
}

type GetAccountHandler struct {
	Command  *GetAccountRequest
	Response GetAccountResponse
	*service.UserService
}

func (h *GetAccountHandler) Handle() error {
	user, err := h.UserService.GetUser(h.Command.ApplicationID, h.Command.Username)
	if user == nil {
		return domainrepository.ErrUserNotFound
	} else if err != nil {
		return err
	}

	if user.ApplicationID != h.Command.ApplicationID || user.Username != h.Command.Username {
		return domainrepository.ErrUserIdentityMismatch
	}

	h.Response = GetAccountResponse{
		&model.UserResponse{
			ID:            user.ID,
			ApplicationID: user.ApplicationID,
			Username:      user.Username,
			Created:       user.Created,
			Updated:       user.Updated,
		},
	}
	return nil
}
