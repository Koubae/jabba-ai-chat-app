package controllers

import (
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/application/user/handlers"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/container"
	domainrepository "github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/domain/user/repository"
	"github.com/Koubae/jabba-ai-chat-app/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AccountControllers struct{}

func (controller *AccountControllers) Get(c *gin.Context) {
	accessToken, exists := c.Get("access_token")
	if !exists {
		c.JSON(400, gin.H{"error": "Access-Token not found"})
		return
	}
	token := accessToken.(*auth.AccessToken)

	request := &handlers.GetAccountRequest{
		ApplicationID: token.ApplicationId,
		Username:      token.Username,
	}

	handler := handlers.GetAccountHandler{Command: request, UserService: container.Container.UserService}

	err := handler.Handle()
	if errors.Is(err, domainrepository.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if errors.Is(err, domainrepository.ErrUserIdentityMismatch) {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, handler.Response.UserResponse)
}
