package controllers

import (
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/application/auth/handlers"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/container"
	domainrepository "github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/domain/user/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Expires     int64  `json:"expires"`
}

func (controller *AuthController) LoginV1(c *gin.Context) {
	var request = handlers.LoginRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	handler := handlers.LoginHandler{Command: request, UserService: container.Container.UserService}

	err := handler.Handle()
	if errors.Is(err, domainrepository.ErrUserNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, handler.Response)
}

func (controller *AuthController) SignUpV1(c *gin.Context) {
	var request = handlers.SignUpRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	handler := handlers.SignUpHandler{Command: request, UserService: container.Container.UserService}
	if err := handler.Handle(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, handler.Response.UserResponse)

}
