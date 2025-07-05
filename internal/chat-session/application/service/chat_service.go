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

	handler := &ChatHandler{
		conn:        conn,
		accessToken: accessToken.AccessToken,
		Session:     session,
		Member:      &model.Member{MemberID: memberID, Channel: channel},
		identity:    identity,
	}

	response, err = handler.Handle(ctx, s.Broadcaster, s.AIBotConnector, s.messageRepository)
	return &response, err
}

const MaxExponentialBackoffRetries = 10
const ExponentialBackoffBaseSleepMs = 100

type ChatHandler struct {
	conn        *websocket.Conn
	accessToken string
	*model.Session
	*model.Member
	identity string

	exponentialBackOffAttempts int
}

func (h *ChatHandler) String() string {
	return h.identity
}

func (h *ChatHandler) Handle(ctx context.Context, broadcaster *bot.Broadcaster, botConnector *bot.AIBotConnector, messageRepository repository.MessageRepository) (string, error) {
	h.sendMessageHistoryToClient(ctx, messageRepository)

	response := "Goodbye"
	var err error
	for {
		messageType, message, err := h.RecV()
		if err != nil {
			break
		}

		// Send the user a message already this is the User's original message!
		h.BroadcastUserPrompt(message, messageRepository, broadcaster, messageType)
		response, err := h.SendPromptToAIBot(broadcaster, botConnector, messageType, string(message))
		if err != nil {
			exponentialBackOffError := h.exponentialBackOff()
			if exponentialBackOffError != nil {
				h.BroadcastSystemMessage(broadcaster, fmt.Sprintf("AI-BOT Is not available at this moment, please try again later"))
				break
			}
			continue
		}
		h.exponentialBackOffAttempts = 0

		reply := response.Reply
		log.Printf("%s (Bot-Reply): %s", h, reply)
		payload := h.createMessagePayload("assistant",
			0, "AI Assistant", "bot", "server", reply)
		payloadBytes, _ := json.Marshal(payload)

		go func() {
			err := messageRepository.AddMessage(context.Background(), h.Session.ApplicationID, h.Session.ID, &payload)
			if err != nil {
				log.Printf("%s Error while adding message to database, message: %+v, error: %s\n", h, message, err)
			}
		}()

		broadcaster.Broadcast(h.Session.ApplicationID, h.Session.ID, messageType, payloadBytes)
	}

	return response, err
}

func (h *ChatHandler) exponentialBackOff() error {
	h.exponentialBackOffAttempts++
	if h.exponentialBackOffAttempts > MaxExponentialBackoffRetries-1 {
		fmt.Printf("%s Exponential Backoff exceeded limit of %d, giving up\n", h, MaxExponentialBackoffRetries)
		return errors.New(fmt.Sprintf("exponential backoff exceeded limit of %d", MaxExponentialBackoffRetries))
	}
	// Simple exponential: 100ms, 200ms, 400ms, 800ms...
	delay := time.Duration(ExponentialBackoffBaseSleepMs*(1<<h.exponentialBackOffAttempts)) * time.Millisecond
	fmt.Printf("Attempt %d failed, retrying in %v...\n", h.exponentialBackOffAttempts+1, delay)
	time.Sleep(delay)
	fmt.Printf("Sleep over...\n")

	return nil
}

func (h *ChatHandler) BroadcastUserPrompt(message []byte, messageRepository repository.MessageRepository, broadcaster *bot.Broadcaster, messageType int) {
	go func() {
		payload := h.createMessagePayload("user",
			h.Member.UserID, h.Member.Username, h.Member.MemberID, h.Member.Channel, string(message))

		go func() {
			err := messageRepository.AddMessage(context.Background(), h.Session.ApplicationID, h.Session.ID, &payload)
			if err != nil {
				log.Printf("%s Error while adding message to database, message: %+v, error: %s\n", h, message, err)
			}
		}()

		payloadBytes, _ := json.Marshal(payload)
		broadcaster.Broadcast(h.Session.ApplicationID, h.Session.ID, messageType, payloadBytes)
	}()
}

func (h *ChatHandler) RecV() (int, []byte, error) {
	messageType, message, err := h.conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
			log.Printf("%s Error reading message | Terminating connection, error: %s\n", h, err)
			err = errors.New("unexpected error while reading message")
		} else {
			log.Printf("%s client closed connection | Terminating connection\n", h)
		}
		return -1, []byte{}, err
	}
	log.Printf("%s received: %s", h.identity, message)
	return messageType, message, err
}

func (h *ChatHandler) sendMessageHistoryToClient(ctx context.Context, messageRepository repository.MessageRepository) {
	messages, err := messageRepository.GetMessages(ctx, h.Session.ApplicationID, h.Session.ID)
	if err != nil {
		log.Printf("%s Error while getting messages from database, error: %s\n", h, err)
		return
	}
	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]
		payloadBytes, _ := json.Marshal(message)
		if err := h.conn.WriteMessage(1, payloadBytes); err != nil {
			fmt.Printf("%s Error sending message to client, message: %+v, error: %s\n", h, message, err)
		}
	}
}

func (h *ChatHandler) SendPromptToAIBot(broadcaster *bot.Broadcaster, botConnector *bot.AIBotConnector, messageType int, message string) (*bot.AIBotResponse, error) {
	response, err := botConnector.SendMessage(context.Background(), h.accessToken, h.Session.ID, message)
	if err == nil {
		return response, nil
	}

	go func() {
		log.Printf("%s Error while calling AI-BOT, error: %s\n", h, err)
		h.BroadcastSystemMessage(broadcaster, fmt.Sprintf("Error while calling AI-BOT, error: %s", err))
	}()
	return nil, err

}

func (h *ChatHandler) createMessagePayload(role string, userID int, username string, memberID string, channel string, message string) model.Message {
	return model.Message{
		ApplicationID: h.Session.ApplicationID,
		SessionID:     h.Session.ID,
		Message:       message,
		Timestamp:     time.Now().Unix(),
		Member: model.Member{
			Role:     role,
			UserID:   userID,
			Username: username,
			MemberID: memberID,
			Channel:  channel,
		},
	}
}

func (h *ChatHandler) BroadcastSystemMessage(broadcaster *bot.Broadcaster, message string) {
	payload := h.createMessagePayload("system",
		0, "System", "system", "server", message)
	payloadBytes, _ := json.Marshal(payload)
	broadcaster.Broadcast(h.Session.ApplicationID, h.Session.ID, 1, payloadBytes)
}
