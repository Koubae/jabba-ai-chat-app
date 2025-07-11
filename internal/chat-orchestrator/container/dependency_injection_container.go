package container

import (
	"context"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/application/service"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/infrastructure/connector"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/infrastructure/database/repository"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/mongodb"
	"log"
)

var Container *DependencyInjectionContainer

func CreateDIContainer() {
	if Container != nil {
		return
	}

	db, err := mongodb.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	applicationRepository := repository.NewApplicationRepository(db)
	userRepository := repository.NewUserRepository(db)
	sessionRepository := repository.NewSessionRepository(db)

	applicationService := service.NewApplicationService(applicationRepository)
	userService := service.NewUserService(userRepository, applicationService)
	sessionService := service.NewSessionService(sessionRepository, applicationService, userService)

	chatSessionConnector := connector.NewChatSessionConnector()

	Container = &DependencyInjectionContainer{
		DB:                    db,
		ApplicationRepository: applicationRepository,
		UserRepository:        userRepository,
		SessionRepository:     sessionRepository,

		ApplicationService:   applicationService,
		UserService:          userService,
		SessionService:       sessionService,
		ChatSessionConnector: chatSessionConnector,
	}
}

func ShutDown() {
	if Container == nil {
		log.Println("DependencyInjectionContainer is not initialized, skipping shutdown")
		return
	}
	Container.Shutdown()
}

type DependencyInjectionContainer struct {
	DB *mongodb.Client
	*repository.ApplicationRepository
	*repository.SessionRepository
	*repository.UserRepository

	*service.ApplicationService
	*service.UserService
	*service.SessionService
	*connector.ChatSessionConnector
}

func (c *DependencyInjectionContainer) Shutdown() {
	log.Println("Shutting down DependencyInjectionContainer and all its resources")

	c.DB.Shutdown(context.Background())
	log.Println("MongoDB database disconnected")

}
