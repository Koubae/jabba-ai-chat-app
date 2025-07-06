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

func NewSessionRepository(db *mongodb.Client) *SessionRepository {
	collection := db.Collection(collections.CollectionSessions)
	return &SessionRepository{db: db, collection: collection}
}

type SessionRepository struct {
	db         *mongodb.Client
	collection *mongo.Collection
}

func (r *SessionRepository) Create(ctx context.Context, session *model.Session) error {
	applicationID, _ := primitive.ObjectIDFromHex(session.ApplicationID)
	userID, _ := primitive.ObjectIDFromHex(session.UserID)

	document := collections.Session{
		ApplicationID: &applicationID,
		UserID:        &userID,
		Name:          session.Name,
	}
	document.OnCreate()
	result, err := r.collection.InsertOne(ctx, document)
	if r.db.IsDuplicateKeyError(err) {
		log.Printf("Session %+v already exists!, error: %s\n", document, err)
		return domainrepository.ErrSessionAlreadyExists
	} else if err != nil {
		log.Printf("Error while creating Sesion %+v, error: %s\n", document, err)
		return domainrepository.ErrSessionOnCreate
	}
	document.OnCreated(result)

	session.ID = document.ID.Hex()
	session.Created = *document.Created
	session.Updated = *document.Updated

	return nil
}

func (r *SessionRepository) GetSession(ctx context.Context, applicationID string, userID string, name string) (*model.Session, error) {
	applicationIDObj, _ := primitive.ObjectIDFromHex(applicationID)
	userIDObj, _ := primitive.ObjectIDFromHex(userID)

	document := &collections.Session{}

	filter := bson.D{
		{"application_id", applicationIDObj},
		{"user_id", userIDObj},
		{"name", name},
	}
	err := r.collection.FindOne(ctx, filter).Decode(document)
	if err != nil {
		log.Printf("Error in GetSession ApplicationID=%s, UserID=%s, Name=%s, error: %s\n", applicationID, userID, name, err)
		return nil, domainrepository.ErrSessionNotFound
	}

	entity := &model.Session{
		ID:            document.ID.Hex(),
		ApplicationID: document.ApplicationID.Hex(),
		UserID:        document.UserID.Hex(),
		Name:          document.Name,
		Created:       *document.Created,
		Updated:       *document.Updated,
	}
	return entity, nil

}

func (r *SessionRepository) GetSessionByID(ctx context.Context, sessionID string) (*model.Session, error) {
	document := &collections.Session{}

	documentID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		log.Printf("Error converting string into MongoDB ObjectID, id %v, erro: %s\n", sessionID, err)
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": documentID}).Decode(document)
	if err != nil {
		log.Printf("Error in GetSessionByID sessionID=%s, error: %s\n", sessionID, err)
		return nil, domainrepository.ErrSessionNotFound
	}

	entity := &model.Session{
		ID:            document.ID.Hex(),
		ApplicationID: document.ApplicationID.Hex(),
		UserID:        document.UserID.Hex(),
		Name:          document.Name,
		Created:       *document.Created,
		Updated:       *document.Updated,
	}
	return entity, nil
}

func (r *SessionRepository) ListWithPagination(
	ctx context.Context,
	applicationID string,
	userID string,
	limit int64,
	offset int64,
) ([]*model.Session, error) {
	findOptions := options.Find()
	findOptions.SetSkip(offset)
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{"_id", 1}})

	applicationIDObjID, _ := primitive.ObjectIDFromHex(applicationID)
	userIDObjID, _ := primitive.ObjectIDFromHex(userID)

	filter := bson.D{
		{"application_id", applicationIDObjID},
		{"user_id", userIDObjID},
	}
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list seessions: %w", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Error while closing cursor (during Applpication ListWithPagination), error: %s\n", err)
		}
	}(cursor, ctx)

	var entities []*model.Session
	if err = cursor.All(ctx, &entities); err != nil {
		return nil, fmt.Errorf("failed to decode sesions: %w", err)
	}
	return entities, nil

}
