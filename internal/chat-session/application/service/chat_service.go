package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/repository"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/infrastructure/bot"
	"github.com/Koubae/jabba-ai-chat-app/pkg/auth"
	"github.com/gorilla/websocket"
	"log"
)

func NewChatService(repository repository.SessionRepository, broadcaster *bot.Broadcaster, botConnector *bot.AIBotConnector) *ChatService {
	return &ChatService{
		repository:     repository,
		Broadcaster:    broadcaster,
		AIBotConnector: botConnector,
	}
}

type ChatService struct {
	repository repository.SessionRepository
	*bot.Broadcaster
	*bot.AIBotConnector
}

func (s *ChatService) CreateConnectionAndStartChat(ctx context.Context, conn *websocket.Conn, sessionID string) (*string, error) {
	accessToken, ok := ctx.Value("access_token").(*auth.AccessToken)
	if !ok {
		return nil, fmt.Errorf("access_token not found, cannot create session")
	}
	identity := fmt.Sprintf("[%s][%s][%s (%d)] (WebSocket)",
		accessToken.ApplicationId, sessionID, accessToken.Username, accessToken.UserId)
	fmt.Printf("Created WebSocket connection %s\n", identity)

	session, _ := s.repository.Get(ctx, accessToken.ApplicationId, sessionID)
	if session == nil {
		log.Printf("Session %+v does not exists for %s\n", session, identity)
		return nil, errors.New("session does not exists, you must create one first")
	}

	s.Broadcaster.Connect(accessToken.ApplicationId, sessionID, accessToken.UserId, accessToken.Username, conn)

	var err error
	response := "Goodbye"

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Printf("%s Error reading message | Terminating connection, error: %s\n", identity, err)
				err = errors.New("unexpected error while reading message")
			} else {
				log.Printf("%s client closed connection | Terminating connection\n", identity)
			}
			break
		}

		log.Printf("%s received: %s", identity, message)
		s.Broadcaster.Broadcast(accessToken.ApplicationId, sessionID, messageType, message)
	}

	return &response, err
}
