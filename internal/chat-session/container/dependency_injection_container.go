package container

import (
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/application/service"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/infrastructure/bot"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/infrastructure/database/repository"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/redis"
	"log"
)

var Container *DependencyInjectionContainer

func CreateDIContainer() {
	if Container != nil {
		return
	}

	db, err := redis.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	aiBotConnector := bot.NewAIBotConnector()
	sessionRepository := repository.NewSessionRepository(db)
	sessionService := service.NewSessionService(sessionRepository, aiBotConnector)

	Container = &DependencyInjectionContainer{
		DB:                db,
		SessionRepository: sessionRepository,
		SessionService:    sessionService,
		AIBotConnector:    aiBotConnector,
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
	DB *redis.Client
	*repository.SessionRepository
	*service.SessionService
	*bot.AIBotConnector
}

func (c *DependencyInjectionContainer) Shutdown() {
	log.Println("Shutting down DependencyInjectionContainer and all its resources")

	c.DB.Shutdown()
	log.Println("Redis database disconnected")

}
