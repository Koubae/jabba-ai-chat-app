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

	broadcaster := bot.NewBroadcaster()
	aiBotConnector := bot.NewAIBotConnector()
	sessionRepository := repository.NewSessionRepository(db)
	messageRepository := repository.NewMessageRepository(db)
	sessionService := service.NewSessionService(sessionRepository, aiBotConnector)

	chatService := service.NewChatService(sessionRepository, messageRepository, broadcaster, aiBotConnector)

	Container = &DependencyInjectionContainer{
		DB:                db,
		SessionRepository: sessionRepository,
		MessageRepository: messageRepository,
		Broadcaster:       broadcaster,
		AIBotConnector:    aiBotConnector,
		SessionService:    sessionService,
		ChatService:       chatService,
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
	*repository.MessageRepository
	*bot.Broadcaster
	*bot.AIBotConnector
	*service.SessionService
	*service.ChatService
}

func (c *DependencyInjectionContainer) Shutdown() {
	log.Println("Shutting down DependencyInjectionContainer and all its resources")

	c.DB.Shutdown()
	log.Println("Redis database disconnected")

}
