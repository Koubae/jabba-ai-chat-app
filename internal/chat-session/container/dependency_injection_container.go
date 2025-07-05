package container

import (
	"context"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/application/service"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/infrastructure/database/repository"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/redis"
	"log"
	"time"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sessionRepository := repository.NewSessionRepository(db, ctx)
	sessionService := service.NewSessionService(sessionRepository)

	Container = &DependencyInjectionContainer{
		DB:                db,
		SessionRepository: sessionRepository,
		SessionService:    sessionService,
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
}

func (c *DependencyInjectionContainer) Shutdown() {
	log.Println("Shutting down DependencyInjectionContainer and all its resources")

	c.DB.Shutdown()
	log.Println("Redis database disconnected")

}
