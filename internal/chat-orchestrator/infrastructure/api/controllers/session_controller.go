package controllers

import (
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/application/handlers"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/container"

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

//func (controller *SessionController) Get(c *gin.Context) {
//	name := c.Param("name")
//
//	request := handlers.GetApplicationRequest{Name: name}
//	handler := handlers.GetApplicationHandler{Command: request, ApplicationService: container.Container.ApplicationService}
//
//	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
//	defer cancel()
//
//	err := handler.Handle(ctx)
//	if errors.Is(err, domainrepository.ErrApplicationNotFound) {
//		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
//		return
//	} else if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	accessToken, _ := c.Get("access_token")
//	accessTokenObj := accessToken.(*auth.AccessToken)
//
//	application := handler.Response
//	if application.Name != accessTokenObj.ApplicationId {
//		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
//		return
//	}
//
//	c.JSON(http.StatusOK, handler.Response)
//
//}
//func (controller *SessionController) List(c *gin.Context) {
//	var request = handlers.ListApplicationRequest{}
//	if err := request.Validate(c); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	handler := handlers.ListApplicationHandler{Command: request, ApplicationService: container.Container.ApplicationService}
//
//	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
//	defer cancel()
//
//	err := handler.Handle(ctx)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, handler.Response)
//
//}
