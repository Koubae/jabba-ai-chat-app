package controllers

import (
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/application/auth/handlers"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/di_container"
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

	container := di_container.Container
	handler := handlers.LoginHandler{Command: request, UserService: container.UserService}

	err := handler.Handle()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	container := di_container.Container
	handler := handlers.SignUpHandler{Command: request, UserService: container.UserService}
	if err := handler.Handle(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, handler.Response.User)

}
