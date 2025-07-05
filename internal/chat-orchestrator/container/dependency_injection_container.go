package container

import (
	"context"
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

	Container = &DependencyInjectionContainer{
		DB: db,
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
}

func (c *DependencyInjectionContainer) Shutdown() {
	log.Println("Shutting down DependencyInjectionContainer and all its resources")

	c.DB.Shutdown(context.Background())
	log.Println("MongoDB database disconnected")

}
