package repository

import (
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/domain/user/model"
	domainrepository "github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/domain/user/repository"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/settings"
	_ "github.com/Koubae/jabba-ai-chat-app/pkg/common/testings"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/utils"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/mysql"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// /////////////////////////
	//			SetUp
	// /////////////////////////
	settings.NewConfig()

	client, err := mysql.NewClient()
	if err != nil {
		panic(err.Error())
	}
	defer client.Shutdown()

	log.Println(client)

	// /////////////////////////
	//			Tests
	// /////////////////////////

	code := m.Run()

	// /////////////////////////
	//			CleanUp
	// /////////////////////////
	// Cleanup after all tests

	_, err = client.DB.Exec("TRUNCATE users")
	if err != nil {
		log.Fatalf("failed to TRUNCATE users database %s: %v", client.Config.DBName, err)
		return
	}
	log.Printf("user table %s truncated successfully\n", client.Config.DBName)

	os.Exit(code)
}

const HashedPassword = "$2a$10$GPeYnQMl9mGX1hvIrqTIjeJmPOESnUHFe39Ksm0HifPU8r9YchbbC"

func TestUserRepository(t *testing.T) {
	client := mysql.GetClient()
	repository := NewUserRepository(client)

	applicationID := "1234567890" + utils.RandomString(5)
	username := "integration-tests" + utils.RandomString(5)
	user := &model.User{
		ApplicationID: applicationID,
		Username:      username,
		PasswordHash:  HashedPassword,
	}

	t.Run("Create", func(t *testing.T) {
		err := repository.Create(user)
		assert.NoError(t, err)
		assert.NotEqual(t, int64(0), user.ID)
	})
	t.Run("Create On Duplicate must return USER_ALREADY_EXISTS error", func(t *testing.T) {
		err := repository.Create(user)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainrepository.ErrUserAlreadyExists)
	})

	t.Run("GetByID", func(t *testing.T) {
		user, err := repository.GetByID(user.ID)

		assert.NoError(t, err)
		assert.Equal(t, applicationID, user.ApplicationID)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, HashedPassword, user.PasswordHash)
	})
	t.Run("GetByID not found", func(t *testing.T) {
		userNotFound, err := repository.GetByID(99999)

		assert.Nil(t, userNotFound)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainrepository.ErrUserNotFound)
	})

	t.Run("GetByUsername", func(t *testing.T) {
		user, err := repository.GetByUsername(applicationID, username)

		assert.NoError(t, err)
		assert.Equal(t, applicationID, user.ApplicationID)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, HashedPassword, user.PasswordHash)
	})
	t.Run("GetByUsername not found", func(t *testing.T) {
		userNotFound, err := repository.GetByUsername("potato", username)

		assert.Nil(t, userNotFound)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainrepository.ErrUserNotFound)

		userNotFound, err = repository.GetByUsername(applicationID, "potato")

		assert.Nil(t, userNotFound)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainrepository.ErrUserNotFound)
	})

}
