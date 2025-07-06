package controllers

import (
	"context"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/container"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type ChatController struct{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: Check how to fix this before go in production!
	},
}

func (controller *ChatController) CreateConnection(c *gin.Context) {
	sessionID := c.Param("session_id")

	memberID := c.Query("member_id")
	channel := c.Query("channel")

	if memberID == "" {
		c.JSON(400, gin.H{"error": "member_id query parameter is required"})
		return
	} else if channel == "" {
		c.JSON(400, gin.H{"error": "channel query parameter is required"})
		return
	}

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

	chatService := container.Container.ChatService
	response, err := chatService.CreateConnectionAndStartChat(ctx, conn, sessionID, memberID, channel)
	log.Printf("Chat connection closed, response: %v, error: %v\n", *response, err)
}
