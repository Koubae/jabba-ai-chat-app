package repository

import (
	"context"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
	domainrepository "github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/repository"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/infrastructure/database/collections"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/settings"
	_ "github.com/Koubae/jabba-ai-chat-app/pkg/common/testings"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/utils"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/mongodb"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// /////////////////////////
	//			SetUp
	// /////////////////////////
	settings.NewConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongodb.NewClient()
	if err != nil {
		panic(err.Error())
	}
	defer client.Shutdown(ctx)

	// Load All collections
	collectionApplications := client.Collection(collections.CollectionApplications)
	if err = client.CreateUniqueIndex(collectionApplications, ctx, "name"); err != nil {
		log.Fatalf("MongoDB error while creating unique index, error %v\n", err)
	}
	// /////////////////////////
	//			Tests
	// /////////////////////////

	code := m.Run()

	// /////////////////////////
	//			CleanUp
	// /////////////////////////
	// Cleanup after all tests

	// Drop All collections
	_ = collectionApplications.Drop(ctx)

	client.Shutdown(ctx)
	os.Exit(code)
}

func TestApplicationRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	db := mongodb.GetClient()
	repository := NewApplicationRepository(db)

	name := "applicationOne-test" + utils.RandomString(10)
	applicationOne := &model.Application{Name: name}
	err := repository.Create(ctx, applicationOne)
	assert.NoError(t, err)

	applicationID := applicationOne.ID

	t.Run("Create", func(t *testing.T) {
		name := "application-test" + utils.RandomString(10)
		application := &model.Application{Name: name}

		err := repository.Create(ctx, application)
		assert.NoError(t, err)
		assert.NotEmpty(t, application.ID)
		assert.NotEmpty(t, application.Updated)
		assert.NotEmpty(t, application.Created)
	})
	t.Run("CreateOnDuplicate", func(t *testing.T) {
		name := "application-test" + utils.RandomString(10)
		application := &model.Application{Name: name}

		err := repository.Create(ctx, application)
		assert.NoError(t, err)
		assert.NotEmpty(t, application.ID)

		application2 := &model.Application{Name: name}
		err2 := repository.Create(ctx, application2)
		assert.Error(t, err2)
		assert.ErrorIs(t, err2, domainrepository.ErrApplicationAlreadyExists)

	})

	t.Run("GetByID", func(t *testing.T) {
		application, err := repository.GetByID(ctx, applicationID)
		assert.NoError(t, err)
		assert.Equal(t, applicationOne.ID, application.ID)
		assert.Equal(t, applicationOne.Name, application.Name)

		t.Run("GetByIDNotFound", func(t *testing.T) {
			objectID := primitive.NewObjectID().Hex()
			application, err := repository.GetByID(ctx, objectID)
			assert.ErrorIs(t, err, domainrepository.ErrApplicationNotFound)
			assert.Nil(t, application)
		})

		t.Run("GetByName", func(t *testing.T) {
			application, err := repository.GetByName(ctx, applicationOne.Name)
			assert.NoError(t, err)
			assert.Equal(t, applicationOne.ID, application.ID)
			assert.Equal(t, applicationOne.Name, application.Name)
		})
		t.Run("GetByNameNotFound", func(t *testing.T) {
			application, err := repository.GetByName(ctx, "potato")
			assert.ErrorIs(t, err, domainrepository.ErrApplicationNotFound)
			assert.Nil(t, application)
		})
	})
}
