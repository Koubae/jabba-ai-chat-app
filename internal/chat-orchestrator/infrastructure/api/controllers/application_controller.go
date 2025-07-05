package controllers

import (
	"context"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/application/handlers"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/container"
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

}
