package controllers

import (
	"context"
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/application/handlers"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/container"
	domainrepository "github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/repository"
	"github.com/Koubae/jabba-ai-chat-app/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ApplicationController struct{}

func (controller *ApplicationController) Create(c *gin.Context) {
	name := c.Param("name")

	request := handlers.CreateApplicationRequest{Name: name}
	handler := handlers.CreateApplicationHandler{Command: request, ApplicationService: container.Container.ApplicationService}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	err := handler.Handle(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, handler.Response)
}

func (controller *ApplicationController) Get(c *gin.Context) {
	name := c.Param("name")

	request := handlers.GetApplicationRequest{Name: name}
	handler := handlers.GetApplicationHandler{Command: request, ApplicationService: container.Container.ApplicationService}

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

	accessToken, _ := c.Get("access_token")
	accessTokenObj := accessToken.(*auth.AccessToken)

	application := handler.Response
	if application.Name != accessTokenObj.ApplicationId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, handler.Response)

}
