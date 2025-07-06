package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/model"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/repository"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/utils"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/redis"
	redisadapter "github.com/redis/go-redis/v9"

	"time"
)

const SessionCacheKey = "_SESSION_CACHE_KEY_"

func NewSessionRepository(db *redis.Client) *SessionRepository {
	cacheServicePrefix := utils.GetEnvString("CACHE_SERVICE_PREFIX", "chat_session:")
	ttlSeconds := utils.GetEnvInt("CACHE_TTL_SECONDS", 1800)

	return &SessionRepository{
		db:                 db,
		cacheServicePrefix: cacheServicePrefix,
		ttlSeconds:         time.Duration(ttlSeconds) * time.Second,
	}
}

type SessionRepository struct {
	db                 *redis.Client
	cacheServicePrefix string
	ttlSeconds         time.Duration
}

func (r *SessionRepository) Create(ctx context.Context, session *model.Session, identityID int64) error {
	document, err := json.Marshal(session)
	if err != nil {
		return err
	}

	key := r.getCacheKey(session.ApplicationID, session.ID, identityID)
	return r.db.DB.Set(ctx, key, document, r.ttlSeconds).Err()
}

func (r *SessionRepository) Get(
	ctx context.Context,
	applicationID string,
	sessionID string,
	identityID int64,
) (*model.Session, error) {
	key := r.getCacheKey(applicationID, sessionID, identityID)

	document, err := r.db.DB.Get(ctx, key).Bytes()
	if errors.Is(err, redisadapter.Nil) || document == nil {
		return nil, repository.ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	var session model.Session
	if err := json.Unmarshal(document, &session); err != nil {
		return nil, repository.ErrSessionParse
	}
	return &session, nil
}

func (r *SessionRepository) getCacheKey(applicationID string, sessionID string, identityID int64) string {
	pattern := "%s%s:%s:%s:%d"
	return fmt.Sprintf(pattern, r.cacheServicePrefix, SessionCacheKey, applicationID, sessionID, identityID)

}
