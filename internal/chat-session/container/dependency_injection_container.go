package container

import (
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

	Container = &DependencyInjectionContainer{
		DB: db,
	}
}

type DependencyInjectionContainer struct {
	DB *redis.Client
}

// TODO add stutdown!
