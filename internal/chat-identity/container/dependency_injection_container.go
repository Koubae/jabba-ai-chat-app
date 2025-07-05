package container

import (
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/application/service"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/infrastructure/database/repository"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/mysql"
	"log"
)

var Container *DependencyInjectionContainer

func CreateDIContainer() {
	if Container != nil {
		return
	}

	db, err := mysql.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)

	Container = &DependencyInjectionContainer{
		DB:             db,
		UserRepository: userRepository,
		UserService:    userService,
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
	DB *mysql.Client
	*repository.UserRepository
	*service.UserService
}

func (c *DependencyInjectionContainer) Shutdown() {
	log.Println("Shutting down DependencyInjectionContainer and all its resources")

	c.DB.Shutdown()
	log.Println("MySQL database disconnected")

}
