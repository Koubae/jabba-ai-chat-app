package main

import (
	"context"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/infrastructure/database/collections"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/settings"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/mongodb"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func init() {
	err := godotenv.Load(".env.chat-orchestrator")
	if err != nil {
		panic(err.Error())
	}
	settings.NewConfig()

	log.SetFlags(log.Ldate | log.Ltime)

}

func main() {
	settings.GetConfig()

	client, err := mongodb.NewClient()
	if err != nil {
		panic(err.Error())
	}
	log.Println(client)

	ctx := context.Background()

	databases, err := client.ListDatabases(ctx)
	if err != nil {
		log.Printf("MongoDB error while listing databases, error %v\n", err)
	}
	log.Printf("MongoDB databases: %v\n", databases)

	// Applications
	collectionApplications := client.Collection(collections.CollectionApplications)
	err = client.CreateUniqueIndex(collectionApplications, ctx, "name")
	if err != nil {
		log.Printf("MongoDB error while creating index for Applications collection, error %v\n", err)
	}

	// Users
	collectionUsers := client.Collection(collections.CollectionUsers)
	err = client.CreateCompoundUniqueIndex(collectionUsers, ctx, []string{"application_id", "username"})
	if err != nil {
		log.Printf("MongoDB error while creating compound unique index of users, error %v\n", err)
	}

	// Sessions
	collectionSessions := client.Collection(collections.CollectionSessions)
	err = client.CreateIndex(collectionSessions, ctx, []string{"application_id"}, []int{1})
	if err != nil {
		log.Printf("MongoDB error while creating Index in sessions collections, error %v\n", err)
	}

	// Members
	collectionMembers := client.Collection(collections.CollectionMembers)
	err = client.CreateIndex(collectionMembers, ctx, []string{"session_id"}, []int{1})
	if err != nil {
		log.Printf("MongoDB error while creating Index in members collections with session_id, error %v\n", err)
	}
	err = client.CreateIndex(collectionMembers, ctx, []string{"user_id"}, []int{1})
	if err != nil {
		log.Printf("MongoDB error while creating Index in members collections with user_id, error %v\n", err)
	}

	// Messages
	collectionMessages := client.Collection(collections.CollectionMessages)
	err = client.CreateIndex(collectionMessages, ctx, []string{"session_id", "user_id"}, []int{1, 1})
	if err != nil {
		log.Printf("MongoDB error while creating Index in messages collections with session_id, user_id, error %v\n", err)
	}

	shutdown(client)
}

func shutdown(client *mongodb.Client) {
	shutDownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Shutdown(shutDownCtx); err != nil {
		log.Fatalf("MongoDB error while shutting Down, error %v\n", err)
	}
	log.Println("MongoDB shutdown completed")
}
