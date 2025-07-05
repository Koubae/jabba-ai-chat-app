package controllers

import (
	"context"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/application/handlers"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/container"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type SessionController struct{}

func (controller *SessionController) CreateSession(c *gin.Context) {
	var request = handlers.CreateSessionRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if accessToken, exists := c.Get("access_token"); exists {
		ctx = context.WithValue(ctx, "access_token", accessToken)
	}

	handler := handlers.CreateSessionHandler{
		Command:        request,
		SessionService: container.Container.SessionService,
	}

	err := handler.Handle(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, handler.Response)
}
