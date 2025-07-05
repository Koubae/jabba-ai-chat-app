package service

import (
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/model"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/repository"
	"log"
	"time"
)

func NewSessionService(repository repository.SessionRepository) *SessionService {
	return &SessionService{
		repository: repository,
	}
}

type SessionService struct {
	repository repository.SessionRepository
}

func (s *SessionService) CreateSession(applicationID string, sessionID string, name string) (*model.Session, error) {
	var session *model.Session

	session = &model.Session{
		ApplicationID: applicationID,
		ID:            sessionID,
		Name:          name,
		Created:       time.Now().UTC(),
		Updated:       time.Now().UTC(),
	}

	err := s.repository.Create(session)
	if err != nil {
		return nil, err
	}

	log.Printf("Session %+v created successfully\n", session)
	return session, nil
}

func (s *SessionService) GetSession(applicationID string, sessionID string) (*model.Session, error) {
	session, err := s.repository.Get(applicationID, sessionID)
	if err != nil {
		return nil, err
	}
	return session, nil
}
