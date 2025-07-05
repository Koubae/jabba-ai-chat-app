package service

import (
	"context"
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/repository"
	"log"
)

func NewSessionService(
	repository repository.SessionRepository,
	applicationService *ApplicationService,
) *SessionService {
	return &SessionService{
		repository:         repository,
		ApplicationService: applicationService,
	}
}

type SessionService struct {
	repository         repository.SessionRepository
	ApplicationService *ApplicationService
}

func (s *SessionService) Create(ctx context.Context, applicationID string, IdentityID int64, name string) (*model.Session, error) {

	userID := fmt.Sprintf("todo-%d", IdentityID)

	application, err := s.ApplicationService.Get(ctx, applicationID)
	if err != nil {
		log.Printf("Application %s not found\n", applicationID)
		return nil, err
	}

	session := &model.Session{
		ApplicationID: application.ID,
		UserID:        userID,
		Name:          name,
	}
	err = s.repository.Create(ctx, session)
	if err != nil {
		log.Printf("Error while creating session %+v, erorr: %s\n", session, err)
		return nil, err
	}

	log.Printf("Session %+v created successfully\n", session)
	return session, nil
}

func (s *SessionService) Get(ctx context.Context, applicationID string, IdentityID int64, name string) (*model.Session, error) {
	userID := fmt.Sprintf("todo-%d", IdentityID)
	application, err := s.ApplicationService.Get(ctx, applicationID)
	if err != nil {
		log.Printf("Application %s not found\n", applicationID)
		return nil, err
	}

	session, err := s.repository.GetSession(ctx, application.ID, userID, name)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) List(
	ctx context.Context,
	applicationID string,
	IdentityID int64,
	limit int64,
	offset int64,
) ([]*model.Session, error) {
	userID := fmt.Sprintf("todo-%d", IdentityID)
	application, err := s.ApplicationService.Get(ctx, applicationID)
	if err != nil {
		log.Printf("Application %s not found\n", applicationID)
		return nil, err
	}

	sessions, err := s.repository.ListWithPagination(ctx, application.ID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	if sessions == nil {
		sessions = []*model.Session{}
	}

	return sessions, nil
}
