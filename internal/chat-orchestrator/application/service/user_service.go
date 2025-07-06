package service

import (
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/repository"
	"golang.org/x/net/context"
	"log"
)

func NewUserService(repository repository.UserRepository, applicationService *ApplicationService) *UserService {
	return &UserService{repository: repository, ApplicationService: applicationService}
}

type UserService struct {
	repository         repository.UserRepository
	ApplicationService *ApplicationService
}

func (s *UserService) Create(
	ctx context.Context,
	applicationID string,
	identityID int64,
	username string,
) (*model.User, error) {
	application, err := s.ApplicationService.Get(ctx, applicationID)
	if err != nil {
		log.Printf("Application %s not found\n", applicationID)
		return nil, err
	}

	user := &model.User{
		ApplicationID: application.ID,
		IdentityID:    identityID,
		Username:      username,
	}
	err = s.repository.Create(ctx, user)
	if err != nil {
		log.Printf("Error while creating user %+v, erorr: %s\n", user, err)
		return nil, err
	}

	log.Printf("User %+v created successfully\n", user)
	return user, nil
}

func (s *UserService) Get(ctx context.Context, applicationID string, IdentityID int64) (*model.User, error) {
	application, err := s.ApplicationService.Get(ctx, applicationID)
	if err != nil {
		log.Printf("Application %s not found\n", applicationID)
		return nil, err
	}

	session, err := s.repository.Get(ctx, application.ID, IdentityID)
	if err != nil {
		return nil, err
	}

	return session, nil
}
