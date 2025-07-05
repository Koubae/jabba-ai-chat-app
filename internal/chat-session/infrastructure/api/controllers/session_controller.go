package controllers

import (
	"context"
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/application/handlers"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/container"
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

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, restrict this!
	},
}

func (controller *SessionController) CreateConnection(c *gin.Context) {
	sessionID := c.Param("session_id")

	fmt.Println(sessionID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if accessToken, exists := c.Get("access_token"); exists {
		ctx = context.WithValue(ctx, "access_token", accessToken)
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("received: %s", message)

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
