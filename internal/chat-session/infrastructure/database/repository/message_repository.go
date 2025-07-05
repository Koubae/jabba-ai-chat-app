package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/model"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/utils"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/redis"

	"time"
)

const MessagesCacheKey = "_MESSAGES_CACHE_KEY_"

func NewMessageRepository(db *redis.Client) *MessageRepository {
	cacheServicePrefix := utils.GetEnvString("CACHE_SERVICE_PREFIX", "chat_session:")
	ttlSeconds := utils.GetEnvInt("CACHE_TTL_SECONDS", 1800)
	messageStorageLength := utils.GetEnvInt("CACHE_MESSAGE_STORAGE_LENGTH", 300)

	return &MessageRepository{
		db:                   db,
		cacheServicePrefix:   cacheServicePrefix,
		ttlSeconds:           time.Duration(ttlSeconds) * time.Second,
		messageStorageLength: int64(messageStorageLength),
	}
}

type MessageRepository struct {
	db                   *redis.Client
	cacheServicePrefix   string
	ttlSeconds           time.Duration
	messageStorageLength int64
}

func (r *MessageRepository) AddMessage(ctx context.Context, applicationID string, sessionID string, message *model.Message) error {
	key := r.getCacheKey(applicationID, sessionID)

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	pipe := r.db.DB.TxPipeline()
	pipe.LPush(ctx, key, messageJSON)               // Add a new message to the front of the list
	pipe.LTrim(ctx, key, 0, r.messageStorageLength) // Keep only the latest 300 messages
	pipe.Expire(ctx, key, r.ttlSeconds)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *MessageRepository) GetMessages(ctx context.Context, applicationID string, sessionID string) ([]*model.Message, error) {
	key := r.getCacheKey(applicationID, sessionID)

	result := r.db.DB.LRange(ctx, key, 0, -1) // Get all messages (newest first)
	if result.Err() != nil {
		return nil, result.Err()
	}

	messages := make([]*model.Message, 0, len(result.Val()))
	for _, messageJSON := range result.Val() {
		var message model.Message
		if err := json.Unmarshal([]byte(messageJSON), &message); err != nil {
			continue // Skip invalid messages
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

func (r *MessageRepository) getCacheKey(applicationID string, sessionID string) string {
	pattern := "%s%s:%s:%s"
	return fmt.Sprintf(pattern, r.cacheServicePrefix, MessagesCacheKey, applicationID, sessionID)

}
