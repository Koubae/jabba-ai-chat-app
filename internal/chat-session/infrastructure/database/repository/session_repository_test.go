package repository

import (
	"context"
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/model"
	domainrepository "github.com/Koubae/jabba-ai-chat-app/internal/chat-session/domain/repository"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/settings"
	_ "github.com/Koubae/jabba-ai-chat-app/pkg/common/testings"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/utils"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/redis"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// /////////////////////////
	//			SetUp
	// /////////////////////////
	settings.NewConfig()

	client, err := redis.NewClient()
	if err != nil {
		panic(err.Error())
	}
	defer client.Shutdown()

	// /////////////////////////
	//			Tests
	// /////////////////////////

	code := m.Run()

	// /////////////////////////
	//			CleanUp
	// /////////////////////////
	// Cleanup after all tests

	os.Exit(code)
}

func TestRedisSessionRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 120*time.Second)
	defer cancel()

	db := redis.GetClient()
	repository := NewSessionRepository(db)

	identityID := int64(1)
	owner := &model.Member{
		Role:     "user",
		UserID:   identityID,
		Username: "test-user",
		MemberID: "device-1234",
		Channel:  "mobile",
	}

	applicationID := "application-id-test" + utils.RandomString(20)
	sessionName := "session-test" + utils.RandomString(20)

	t.Run(
		"Create", func(t *testing.T) {
			session := &model.Session{
				ApplicationID: applicationID,
				ID:            fmt.Sprintf("session-id-test-%d", utils.RandInt(1, 99999)),
				Name:          sessionName,
				Owner:         owner,
				Created:       time.Now().UTC(),
				Updated:       time.Now().UTC(),
			}

			err := repository.Create(ctx, session, identityID)
			assert.NoError(t, err)
		},
	)

	t.Run(
		"Get", func(t *testing.T) {
			session := &model.Session{
				ApplicationID: applicationID,
				ID:            fmt.Sprintf("session-id-test-%d", utils.RandInt(1, 99999)),
				Name:          sessionName,
				Owner:         owner,
				Created:       time.Now().UTC(),
				Updated:       time.Now().UTC(),
			}

			err := repository.Create(ctx, session, identityID)
			assert.NoError(t, err)

			sessionInCache, err := repository.Get(ctx, session.ApplicationID, session.ID, identityID)
			assert.NoError(t, err)
			assert.Equal(t, session.ID, sessionInCache.ID)
			assert.Equal(t, session.ApplicationID, sessionInCache.ApplicationID)
			assert.Equal(t, session.Name, sessionInCache.Name)
			assert.Equal(t, session.Created, sessionInCache.Created)
			assert.Equal(t, session.Updated, sessionInCache.Updated)

		},
	)

	t.Run(
		"GetNotFound", func(t *testing.T) {
			sessionInCache, err := repository.Get(ctx, applicationID, "potato", identityID)
			assert.ErrorIs(t, err, domainrepository.ErrSessionNotFound)
			assert.Nil(t, sessionInCache)

		},
	)

}
