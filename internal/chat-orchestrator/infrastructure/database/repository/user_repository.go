package repository

import (
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/model"
	domainrepository "github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/domain/application/repository"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-orchestrator/infrastructure/database/collections"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"log"
)

func NewUserRepository(db *mongodb.Client) *UserRepository {
	collection := db.Collection(collections.CollectionUsers)
	return &UserRepository{db: db, collection: collection}
}

type UserRepository struct {
	db         *mongodb.Client
	collection *mongo.Collection
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	applicationID, _ := primitive.ObjectIDFromHex(user.ApplicationID)

	document := collections.User{
		ApplicationID: &applicationID,
		IdentityID:    user.IdentityID,
		Username:      user.Username,
	}
	document.OnCreate()
	result, err := r.collection.InsertOne(ctx, document)
	if r.db.IsDuplicateKeyError(err) {
		log.Printf("User %+v already exists!, error: %s\n", document, err)
		return domainrepository.ErrUserAlreadyExists
	} else if err != nil {
		log.Printf("Error while creating User %+v, error: %s\n", document, err)
		return domainrepository.ErrUserOnCreate
	}
	document.OnCreated(result)

	user.ID = document.ID.Hex()
	user.Created = *document.Created
	user.Updated = *document.Updated

	return nil
}

func (r *UserRepository) Get(ctx context.Context, applicationID string, IdentityID int64) (*model.User, error) {
	applicationIDObj, _ := primitive.ObjectIDFromHex(applicationID)

	document := &collections.User{}
	filter := bson.D{
		{"application_id", applicationIDObj},
		{"identity_id", IdentityID},
	}
	err := r.collection.FindOne(ctx, filter).Decode(document)
	if err != nil {
		log.Printf("Error in Get User ApplicationID=%s, IdentityID=%d, error: %s\n", applicationID, IdentityID, err)
		return nil, domainrepository.ErrUserNotFound
	}

	entity := &model.User{
		ID:            document.ID.Hex(),
		ApplicationID: document.ApplicationID.Hex(),
		IdentityID:    document.IdentityID,
		Username:      document.Username,
		Created:       *document.Created,
		Updated:       *document.Updated,
	}
	return entity, nil
}
