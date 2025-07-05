package repository

import (
	"context"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
	domainrepository "github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/repository"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/infrastructure/database/collections"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type ApplicationRepository struct {
	db         *mongodb.Client
	ctx        context.Context
	collection *mongo.Collection
}

func NewApplicationRepository(db *mongodb.Client, ctx context.Context) *ApplicationRepository {
	collection := db.Collection(collections.CollectionApplications)
	return &ApplicationRepository{db: db, ctx: ctx, collection: collection}
}

func (r *ApplicationRepository) Create(application *model.Application) error {
	document := collections.Application{
		Name: application.Name,
	}
	document.OnCreate()

	result, err := r.collection.InsertOne(r.ctx, application)
	if err != nil {
		log.Printf("Error while creating application %+v, error: %s\n", document, err)
		return domainrepository.ErrApplicationOnCreate
	}
	document.OnCreated(result)

	application.ID = document.ID.Hex()
	application.Created = *document.Created
	application.Updated = *document.Updated

	return nil

}

func (r *ApplicationRepository) GetByID(id string) (*model.Application, error) {
	document := &collections.Application{}

	documentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error converting string into MongoDB ObjectID, id %v, erro: %s\n", id, err)
		return nil, err
	}

	err = r.collection.FindOne(r.ctx, bson.M{"_id": documentID}).Decode(document)
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

func (r *ApplicationRepository) GetByName(name string) (*model.Application, error) {
	document := &collections.Application{}
	err := r.collection.FindOne(r.ctx, bson.M{"name": name}).Decode(document)
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
