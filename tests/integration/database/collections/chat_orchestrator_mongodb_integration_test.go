package collections

import (
	"context"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/chat/model"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/infrastructure/database/collections"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/settings"
	_ "github.com/Koubae/jabba-ai-chat-app/pkg/common/testings"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/utils"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/mongodb"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	client, err := mongodb.NewClient()
	if err != nil {
		panic(err.Error())
	}
	log.Println(client)

	// Load All collections
	collectionApplications := client.Collection(model.CollectionApplications)
	collectionUsers := client.Collection(model.CollectionUsers)
	collectionSessions := client.Collection(model.CollectionSessions)
	collectionMembers := client.Collection(model.CollectionMembers)
	collectionMessages := client.Collection(model.CollectionMessages)

	// /////////////////////////
	//			Tests
	// /////////////////////////

	code := m.Run()

	// /////////////////////////
	//			CleanUp
	// /////////////////////////
	// Cleanup after all tests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Drop All collections
	_ = collectionApplications.Drop(ctx)
	_ = collectionUsers.Drop(ctx)
	_ = collectionSessions.Drop(ctx)
	_ = collectionMembers.Drop(ctx)
	_ = collectionMessages.Drop(ctx)

	if err := client.Shutdown(ctx); err != nil {
		log.Fatalf("MongoDB error while shutting Down, error %v\n", err)
	}
	log.Println("MongoDB shutdown completed")

	os.Exit(code)
}

func TestIntegrationMongoDBApplicationCollection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	db := mongodb.GetClient()
	collection := db.Collection(model.CollectionApplications)

	err := db.CreateUniqueIndex(collection, ctx, "name")
	assert.NoError(t, err)

	name := "application-test" + utils.RandomString(5)
	application := &collections.Application{Name: name}

	t.Run("nil_values_expected_when_not_saved", func(t *testing.T) {
		assert.Nil(t, application.ID)
		assert.Nil(t, application.Created)
		assert.Nil(t, application.Updated)
	})

	t.Run("Create", func(t *testing.T) {
		application.OnCreate()
		result, err := collection.InsertOne(ctx, application)
		assert.NoError(t, err)

		application.OnCreated(result)

		expectedID := result.InsertedID
		assert.Equal(t, expectedID, *application.ID)
	})

	t.Run("CreateOnDuplicate", func(t *testing.T) {
		application.OnCreate()
		result, err := collection.InsertOne(ctx, application)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Find", func(t *testing.T) {
		found := &collections.Application{}

		err := collection.FindOne(ctx, bson.M{"_id": application.ID}).Decode(found)
		assert.NoError(t, err)
		assert.Equal(t, application.ID, found.ID)
		assert.Equal(t, application.Name, found.Name)
		assert.WithinDuration(t, *application.Created, *found.Created, time.Second*2)
	})

	t.Run("FindNotFound", func(t *testing.T) {
		objectID := primitive.NewObjectID()
		found := &collections.Application{}
		err := collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(found)

		assert.Error(t, err, mongo.ErrNoDocuments)
		assert.Nil(t, found.ID)
	})

	t.Run("Update", func(t *testing.T) {
		// Create a new Application
		nameUpdate := "application-test" + utils.RandomString(5)
		application = &collections.Application{Name: nameUpdate}

		application.OnCreate()
		result, err := collection.InsertOne(ctx, application)
		assert.NoError(t, err)

		application.OnCreated(result)
		// ----

		nameNew := application.Name + "-updated"
		application.Name = nameNew
		application.OnUpdate()

		filter := bson.M{"_id": application.ID}
		update := bson.M{
			"$set": bson.M{
				"name":    application.Name,
				"updated": application.Updated,
			},
		}

		result2, err2 := collection.UpdateOne(ctx, filter, update)
		assert.NoError(t, err2)
		assert.Equal(t, int64(1), result2.ModifiedCount)

		found := &collections.Application{}
		err = collection.FindOne(ctx, bson.M{"_id": application.ID}).Decode(found)
		assert.NoError(t, err)
		assert.Equal(t, nameNew, found.Name)
	})

	t.Run("UpdateNotFound", func(t *testing.T) {
		idNotExists := primitive.NewObjectID()
		filter := bson.M{"_id": idNotExists}
		update := bson.M{
			"$set": bson.M{
				"name": "dont-exists",
			},
		}

		result, err := collection.UpdateOne(ctx, filter, update)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), result.ModifiedCount)
	})

	t.Run("UpsertNewObject", func(t *testing.T) {
		nameNew := "application-test-upsert" + utils.RandomString(5)

		filter := bson.M{"name": nameNew}
		update := bson.M{
			"$set": bson.M{
				"name":    nameNew,
				"updated": time.Now().UTC(),
			},
			"$setOnInsert": bson.M{
				"created": time.Now().UTC(),
			},
		}
		opts := options.Update().SetUpsert(true)
		result, err := collection.UpdateOne(ctx, filter, update, opts)
		assert.NoError(t, err)

		id := result.UpsertedID.(primitive.ObjectID)

		found := &collections.Application{}
		err = collection.FindOne(ctx, bson.M{"name": nameNew}).Decode(found)

		assert.NoError(t, err)
		assert.Equal(t, id, *found.ID)
		assert.Equal(t, nameNew, found.Name)
	})

	t.Run("FindOneAndUpdate", func(t *testing.T) {
		nameNew := "application-test-upsert" + utils.RandomString(10)

		filter := bson.M{"name": nameNew}
		update := bson.M{
			"$set": bson.M{
				"name":    nameNew,
				"updated": time.Now().UTC(),
			},
			"$setOnInsert": bson.M{
				"created": time.Now().UTC(),
			},
		}

		var result collections.Application
		opts := options.FindOneAndUpdate().
			SetUpsert(true).
			SetReturnDocument(options.After)
		err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
		assert.NoError(t, err)

		id := result.ID

		found := &collections.Application{}
		err = collection.FindOne(ctx, bson.M{"name": nameNew}).Decode(found)

		assert.NoError(t, err)
		assert.Equal(t, id, found.ID)
		assert.Equal(t, nameNew, found.Name)
	})

	t.Run("Delete", func(t *testing.T) {
		// Create a new Application
		nameUpdate := "application-test" + utils.RandomString(5)
		application = &collections.Application{Name: nameUpdate}

		application.OnCreate()
		result, err := collection.InsertOne(ctx, application)
		assert.NoError(t, err)

		application.OnCreated(result)
		// ----

		filter := bson.M{"_id": application.ID}
		result2, err2 := collection.DeleteOne(ctx, filter)
		assert.NoError(t, err2)
		assert.Equal(t, int64(1), result2.DeletedCount)

		found := &collections.Application{}
		err3 := collection.FindOne(ctx, bson.M{"_id": application.ID}).Decode(found)
		assert.Error(t, err3, mongo.ErrNoDocuments)
		assert.Nil(t, found.ID)
	})

}
