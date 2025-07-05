package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/model"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/repository"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/infrastructure/bot"
	"github.com/Koubae/jabba-ai-chat-app/pkg/auth"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func NewChatService(
	sessionRepository repository.SessionRepository,
	messageRepository repository.MessageRepository,
	broadcaster *bot.Broadcaster,
	botConnector *bot.AIBotConnector,
) *ChatService {
	return &ChatService{
		sessionRepository: sessionRepository,
		messageRepository: messageRepository,
		Broadcaster:       broadcaster,
		AIBotConnector:    botConnector,
	}
}

type ChatService struct {
	sessionRepository repository.SessionRepository
	messageRepository repository.MessageRepository
	*bot.Broadcaster
	*bot.AIBotConnector
}

func (s *ChatService) CreateConnectionAndStartChat(
	ctx context.Context,
	conn *websocket.Conn,
	sessionID string,
	memberID string,
	channel string,
) (*string, error) {
	accessToken, ok := ctx.Value("access_token").(*auth.AccessToken)
	if !ok {
		return nil, fmt.Errorf("access_token not found, cannot create session")
	}
	identity := fmt.Sprintf("[%s][%s]> Username=%s UserID=%d, Member=%s, Channel=%s (WebSocket)",
		accessToken.ApplicationId, sessionID, accessToken.Username, accessToken.UserId, memberID, channel)
	fmt.Printf("Created WebSocket connection %s\n", identity)

	session, _ := s.sessionRepository.Get(ctx, accessToken.ApplicationId, sessionID)
	if session == nil {
		log.Printf("Session %+v does not exists for %s\n", session, identity)
		return nil, errors.New("session does not exists, you must create one first")
	}

	response := "Goodbye"
	var err error
	err = s.Broadcaster.Connect(conn, accessToken.ApplicationId, sessionID, accessToken.UserId, accessToken.Username, memberID, channel)
	if err != nil {
		response = "Not Allow to connect"
		log.Printf("%s Could not connect client to broadcaster, will close connection, error: %s\n", identity, err)
		return &response, err
	}
	defer s.Broadcaster.Disconnect(conn, accessToken.ApplicationId, sessionID)

	s.sendMessageHistoryToClient(ctx, conn, accessToken, sessionID, identity)

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

		// Send the user a message already this is the User's original message!
		go func() {
			payload := createMessagePayload(accessToken.ApplicationId, sessionID, "user",
				int(accessToken.UserId), accessToken.Username, memberID, channel, string(message))

			go func() {
				err := s.messageRepository.AddMessage(context.Background(), accessToken.ApplicationId, sessionID, &payload)
				if err != nil {
					log.Printf("%s Error while adding message to database, message: %+v, error: %s\n", identity, message, err)
				}
			}()

			payloadBytes, _ := json.Marshal(payload)
			s.Broadcaster.Broadcast(accessToken.ApplicationId, sessionID, messageType, payloadBytes)
		}()

		response, err := s.AIBotConnector.SendMessage(context.Background(), accessToken.AccessToken, session.ID, string(message))
		if err != nil {
			log.Printf("%s Error while calling AI-BOT, error: %s\n", identity, err)

			payload := createMessagePayload(accessToken.ApplicationId, sessionID, "system",
				0, "Error", "system", "server", fmt.Sprintf("Error while calling AI-BOT, error: %s", err))
			payloadBytes, _ := json.Marshal(payload)
			s.Broadcaster.Broadcast(accessToken.ApplicationId, sessionID, messageType, payloadBytes)
			continue
		}

		reply := response.Reply
		log.Printf("%s (Bot-Reply): %s", identity, reply)
		payload := createMessagePayload(accessToken.ApplicationId, sessionID, "assistant",
			0, "AI Assistant", "bot", "server", reply)
		payloadBytes, _ := json.Marshal(payload)

		go func() {
			err := s.messageRepository.AddMessage(context.Background(), accessToken.ApplicationId, sessionID, &payload)
			if err != nil {
				log.Printf("%s Error while adding message to database, message: %+v, error: %s\n", identity, message, err)
			}
		}()

		s.Broadcaster.Broadcast(accessToken.ApplicationId, sessionID, messageType, payloadBytes)
	}

	return &response, err
}

func createMessagePayload(applicationID string, sessionID string, role string, userID int, username string, memberID string, channel string, message string) model.Message {
	return model.Message{
		ApplicationID: applicationID,
		SessionID:     sessionID,
		Role:          role,
		UserID:        userID,
		Username:      username,
		Message:       message,
		Timestamp:     time.Now().Unix(),
		Member: model.Member{
			MemberID: memberID,
			Channel:  channel,
		},
	}
}

func (s *ChatService) sendMessageHistoryToClient(ctx context.Context, conn *websocket.Conn, accessToken *auth.AccessToken, sessionID string, identity string) {
	messages, err := s.messageRepository.GetMessages(ctx, accessToken.ApplicationId, sessionID)
	if err != nil {
		log.Printf("%s Error while getting messages from database, error: %s\n", identity, err)
		return
	}
	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]
		payloadBytes, _ := json.Marshal(message)
		if err := conn.WriteMessage(1, payloadBytes); err != nil {
			fmt.Printf("%s Error sending message to client, message: %+v, error: %s\n", identity, message, err)
		}
	}
}
