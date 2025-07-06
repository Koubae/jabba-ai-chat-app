package service

import (
	"context"
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/model"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/repository"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/infrastructure/bot"
	"github.com/Koubae/jabba-ai-chat-app/pkg/auth"
	"log"
	"time"
)

func NewSessionService(repository repository.SessionRepository, botConnector *bot.AIBotConnector) *SessionService {
	return &SessionService{
		repository:     repository,
		AIBotConnector: botConnector,
	}
}

type SessionService struct {
	repository repository.SessionRepository
	*bot.AIBotConnector
}

func (s *SessionService) CreateSession(
	ctx context.Context,
	sessionID string,
	name string,
	memberID string,
	channel string,
) (*model.Session, error) {
	accessToken, ok := ctx.Value("access_token").(*auth.AccessToken)
	if !ok {
		return nil, fmt.Errorf("access_token not found, cannot create session")
	}

	var session *model.Session
	session = &model.Session{
		ApplicationID: accessToken.ApplicationId,
		ID:            sessionID,
		Name:          name,
		Owner: &model.Member{
			Role:     "user",
			UserID:   accessToken.UserId,
			Username: accessToken.Username,
			MemberID: memberID,
			Channel:  channel,
		},
		Created: time.Now().UTC(),
		Updated: time.Now().UTC(),
	}

	sessionInCache, _ := s.repository.Get(ctx, session.ApplicationID, session.ID, accessToken.UserId)
	if sessionInCache != nil {
		if !sessionInCache.IsSameOwner(session) {
			log.Printf(
				"Session already exists but belongs to a different owner.\n-Request: %+v\n-Found: %+v\n",
				session,
				sessionInCache,
			)
			return nil, model.ErrIsNotOwnerOfSession
		}

		log.Printf("Session %+v already exists in cache, returning it\n", session)
		return sessionInCache, nil
	}

	go func() {
		response, err := s.AIBotConnector.Hello(context.Background(), accessToken.AccessToken, session.ID)
		if err != nil {
			log.Printf("Error while calling AI-BOT, error: %s\n", err)
		} else {
			log.Printf("Hello Response from AI-BOT (session-created): %+v\n", response)
		}
	}()

	err := s.repository.Create(ctx, session, accessToken.UserId)
	if err != nil {
		log.Printf("Error while creating Session %+v, erorr: %s\n", session, err)
		return nil, err
	}

	log.Printf("Session %+v created successfully\n", session)
	return session, nil
}

func (s *SessionService) GetSession(ctx context.Context, applicationID string, sessionID string, identityID int64) (
	*model.Session,
	error,
) {
	session, err := s.repository.Get(ctx, applicationID, sessionID, identityID)
	if err != nil {
		return nil, err
	}
	return session, nil
}
