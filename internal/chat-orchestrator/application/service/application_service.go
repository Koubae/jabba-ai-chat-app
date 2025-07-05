package service

import (
	"context"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/repository"
	"log"
)

func NewApplicationService(repository repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{
		repository: repository,
	}
}

type ApplicationService struct {
	repository repository.ApplicationRepository
}

func (s *ApplicationService) Create(ctx context.Context, name string) (*model.Application, error) {
	application := &model.Application{Name: name}
	err := s.repository.Create(ctx, application)
	if err != nil {
		log.Printf("Error while creating Application %+v, erorr: %s\n", application, err)
		return nil, err
	}

	log.Printf("Application %+v created successfully\n", application)
	return application, nil
}

func (s *ApplicationService) Get(ctx context.Context, name string) (*model.Application, error) {
	application, err := s.repository.GetByName(ctx, name)
	if err != nil {
		log.Printf("Application %s not found\n", name)
		return nil, err
	}

	return application, nil
}
