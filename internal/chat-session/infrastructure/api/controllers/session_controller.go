package controllers

import (
	"context"
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/application/handlers"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/container"
	"github.com/Koubae/jabba-ai-chat-app/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: Check how to fix this before go in production!
	},
}

func (controller *SessionController) CreateConnection(c *gin.Context) {
	sessionID := c.Param("session_id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	accessToken, exists := c.Get("access_token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}
	ctx = context.WithValue(ctx, "access_token", accessToken)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error while closing WebSocket connection, error: %s\n\n", err)
		}
	}(conn)

	accessTokenObj, _ := ctx.Value("access_token").(*auth.AccessToken)

	identity := fmt.Sprintf("[%s][%s][%s (%d)] (WebSocket)",
		accessTokenObj.ApplicationId, sessionID, accessTokenObj.Username, accessTokenObj.UserId)
	fmt.Printf("Created WebSocket connection %s\n", identity)

	broadcaster := container.Container.Broadcaster
	broadcaster.Connect(accessTokenObj.ApplicationId, sessionID, accessTokenObj.UserId, accessTokenObj.Username, conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Printf("%s Error reading message | Terminating connection, error: %s\n", identity, err)
			} else {
				log.Printf("%s client closed connection | Terminagin connection\n", identity)
			}
			break
		}

		log.Printf("%s received: %s", identity, message)
		broadcaster.Broadcast(accessTokenObj.ApplicationId, sessionID, messageType, message)
	}
}
