package service

import (
	"context"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/repository"
	"log"
)

func NewSessionService(
	repository repository.SessionRepository,
	applicationService *ApplicationService,
	userService *UserService,
) *SessionService {
	return &SessionService{
		repository:         repository,
		ApplicationService: applicationService,
		UserService:        userService,
	}
}

type SessionService struct {
	repository         repository.SessionRepository
	ApplicationService *ApplicationService
	UserService        *UserService
}

func (s *SessionService) Create(ctx context.Context, applicationID string, IdentityID int64, username string, name string) (*model.Session, error) {
	application, err := s.ApplicationService.Get(ctx, applicationID)
	if err != nil {
		log.Printf("Application %s not found\n", applicationID)
		return nil, err
	}

	user, err := s.getOrCreateUser(ctx, applicationID, IdentityID, username, err)
	if err != nil || user == nil {
		log.Printf("User %v in application %s not found even after creating it\n", IdentityID, applicationID)
		return nil, err
	}

	userID := user.ID

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

func (s *SessionService) getOrCreateUser(ctx context.Context, applicationID string, IdentityID int64, username string, err error) (*model.User, error) {
	user, _ := s.UserService.Get(ctx, applicationID, IdentityID)
	if user == nil {
		_, err := s.UserService.Create(ctx, applicationID, IdentityID, username)
		user, err := s.UserService.Get(ctx, applicationID, IdentityID)
		if err != nil || user == nil {
			log.Printf("User %v in application %s not found even after creating it\n", IdentityID, applicationID)
			return nil, err
		}
	}
	if user == nil {
		log.Printf("User %v in application %s not found even after creating it\n", IdentityID, applicationID)
		return nil, err
	}
	return user, nil
}

func (s *SessionService) Get(ctx context.Context, applicationID string, IdentityID int64, name string) (*model.Session, error) {
	application, err := s.ApplicationService.Get(ctx, applicationID)
	if err != nil {
		log.Printf("Application %s not found\n", applicationID)
		return nil, err
	}
	user, err := s.UserService.Get(ctx, applicationID, IdentityID)
	if err != nil || user == nil {
		log.Printf("User %v in application %s not found\n", IdentityID, applicationID)
		return nil, err
	}
	userID := user.ID

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
	application, err := s.ApplicationService.Get(ctx, applicationID)
	if err != nil {
		log.Printf("Application %s not found\n", applicationID)
		return nil, err
	}
	user, err := s.UserService.Get(ctx, applicationID, IdentityID)
	if err != nil || user == nil {
		log.Printf("User %v in application %s not found\n", IdentityID, applicationID)
		return nil, err
	}
	userID := user.ID

	sessions, err := s.repository.ListWithPagination(ctx, application.ID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	if sessions == nil {
		sessions = []*model.Session{}
	}

	return sessions, nil
}
