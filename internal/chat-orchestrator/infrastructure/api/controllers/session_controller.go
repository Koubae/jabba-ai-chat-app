package controllers

import (
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/application/handlers"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/container"
	domainrepository "github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/repository"

	"github.com/Koubae/jabba-ai-chat-app/pkg/auth"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type SessionController struct{}

func (controller *SessionController) Create(c *gin.Context) {
	accessToken, _ := c.Get("access_token")
	accessTokenObj := accessToken.(*auth.AccessToken)
	name := c.Param("name")

	request := handlers.CreateSessionRequest{
		ApplicationID: accessTokenObj.ApplicationId,
		IdentityID:    accessTokenObj.UserId,
		Username:      accessTokenObj.Username,
		Name:          name,
	}
	handler := handlers.CreateSessionHandler{Command: request, SessionService: container.Container.SessionService}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	err := handler.Handle(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, handler.Response)
}

func (controller *SessionController) Get(c *gin.Context) {
	accessToken, _ := c.Get("access_token")
	accessTokenObj := accessToken.(*auth.AccessToken)
	name := c.Param("name")

	request := handlers.GetSessionRequest{
		ApplicationID: accessTokenObj.ApplicationId,
		IdentityID:    accessTokenObj.UserId,
		Name:          name,
	}
	handler := handlers.GetSessionHandler{Command: request, SessionService: container.Container.SessionService}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err := handler.Handle(ctx)
	if errors.Is(err, domainrepository.ErrApplicationNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, handler.Response)

}

func (controller *SessionController) List(c *gin.Context) {
	accessToken, _ := c.Get("access_token")
	accessTokenObj := accessToken.(*auth.AccessToken)

	var request = handlers.ListSessionRequest{
		ApplicationID: accessTokenObj.ApplicationId,
		IdentityID:    accessTokenObj.UserId,
	}
	if err := request.Validate(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	handler := handlers.ListSessionHandler{Command: request, SessionService: container.Container.SessionService}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err := handler.Handle(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, handler.Response)

}

func (controller *SessionController) StartSession(c *gin.Context) {
	var request = handlers.StartSessionRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if accessToken, exists := c.Get("access_token"); exists {
		ctx = context.WithValue(ctx, "access_token", accessToken)
	}

	handler := handlers.StartSessionHandler{
		Command:              request,
		SessionService:       container.Container.SessionService,
		ChatSessionConnector: container.Container.ChatSessionConnector,
	}
	err := handler.Handle(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, handler.Response)
}
