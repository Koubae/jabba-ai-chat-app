package repository

import (
	"context"
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
	domainrepository "github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/repository"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/infrastructure/database/collections"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type ApplicationRepository struct {
	db         *mongodb.Client
	collection *mongo.Collection
}

func NewApplicationRepository(db *mongodb.Client) *ApplicationRepository {
	collection := db.Collection(collections.CollectionApplications)
	return &ApplicationRepository{db: db, collection: collection}
}

func (r *ApplicationRepository) Create(ctx context.Context, application *model.Application) error {
	document := collections.Application{
		Name: application.Name,
	}
	document.OnCreate()

	result, err := r.collection.InsertOne(ctx, application)
	if r.db.IsDuplicateKeyError(err) {
		log.Printf("Application %+v already exists!, error: %s\n", document, err)
		return domainrepository.ErrApplicationAlreadyExists
	} else if err != nil {
		log.Printf("Error while creating application %+v, error: %s\n", document, err)
		return domainrepository.ErrApplicationOnCreate
	}
	document.OnCreated(result)

	application.ID = document.ID.Hex()
	application.Created = *document.Created
	application.Updated = *document.Updated

	return nil

}

func (r *ApplicationRepository) GetByID(ctx context.Context, id string) (*model.Application, error) {
	document := &collections.Application{}

	documentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error converting string into MongoDB ObjectID, id %v, erro: %s\n", id, err)
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": documentID}).Decode(document)
	if err != nil {
		log.Printf("Error in GetByID with id %v, error: %s\n", id, err)
		return nil, domainrepository.ErrApplicationNotFound
	}

	entity := &model.Application{
		ID:      document.ID.Hex(),
		Name:    document.Name,
		Created: *document.Created,
		Updated: *document.Updated,
	}
	return entity, nil

}

func (r *ApplicationRepository) GetByName(ctx context.Context, name string) (*model.Application, error) {
	document := &collections.Application{}
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(document)
	if err != nil {
		log.Printf("Error in GetByID with name %v, error: %s\n", name, err)
		return nil, domainrepository.ErrApplicationNotFound
	}

	entity := &model.Application{
		ID:      document.ID.Hex(),
		Name:    document.Name,
		Created: *document.Created,
		Updated: *document.Updated,
	}
	return entity, nil

}

func (r *ApplicationRepository) ListWithPagination(ctx context.Context, limit, offset int64) ([]*model.Application, error) {
	findOptions := options.Find()
	findOptions.SetSkip(offset)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{"_id", 1}})
	cursor, err := r.collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Error while closing cursor (during Applpication ListWithPagination), error: %s\n", err)
		}
	}(cursor, ctx)

	var applications []*model.Application
	if err = cursor.All(ctx, &applications); err != nil {
		return nil, fmt.Errorf("failed to decode applications: %w", err)
	}
	return applications, nil

}
