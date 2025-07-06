package collections

import (
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionApplications = "applications"
const CollectionUsers = "users"
const CollectionSessions = "sessions"
const CollectionMembers = "members"
const CollectionMessages = "messages"

type Application struct {
	mongodb.EntityID   `bson:",inline"`
	mongodb.Timestamps `bson:",inline"`
	Name               string `bson:"name"`
}

type User struct {
	mongodb.EntityID   `bson:",inline"`
	mongodb.Timestamps `bson:",inline"`
	ApplicationID      *primitive.ObjectID `bson:"application_id"`
	IdentityID         int64               `bson:"identity_id"`
	Username           string              `bson:"username"`
}

type Session struct {
	mongodb.EntityID   `bson:",inline"`
	mongodb.Timestamps `bson:",inline"`
	ApplicationID      *primitive.ObjectID `bson:"application_id"`
	UserID             *primitive.ObjectID `bson:"user_id"`
	Name               string              `bson:"name"`
}

type Member struct {
	mongodb.EntityID   `bson:",inline"`
	mongodb.Timestamps `bson:",inline"`
	SessionID          *primitive.ObjectID `bson:"session_id"`
	UserID             *primitive.ObjectID `bson:"user_id"`
	Channel            string              `bson:"channel"`
}

type Message struct {
	mongodb.EntityID   `bson:",inline"`
	mongodb.Timestamps `bson:",inline"`
	SessionID          *primitive.ObjectID `bson:"session_id"`
	UserID             *primitive.ObjectID `bson:"user_id"`
	Role               string              `bson:"role"`
	Body               string              `bson:"body"`
}
