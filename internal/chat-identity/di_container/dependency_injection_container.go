package di_container

import (
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/application/service"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/infrastructure/database/repository"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/mysql"
	"log"
)

type DependencyInjectionContainer struct {
	DB *mysql.Client
	*repository.UserRepository
	*service.UserService
}

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
