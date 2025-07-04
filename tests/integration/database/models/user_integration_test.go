package models

import (
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/domain/user/model"
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

func TestUserCrud(t *testing.T) {
	client := mysql.GetClient()

	username := "integration-tests" + utils.RandomString(5)
	user := &model.User{ApplicationID: "application-integration-tests", Username: username, PasswordHash: HashedPassword}
	var id int64

	t.Run("Create", func(t *testing.T) {
		query := `
			INSERT INTO users (application_id, username, password_hash) 
			VALUES (?, ?, ?)
		`
		result, err := client.DB.Exec(query, user.ApplicationID, user.Username, user.PasswordHash)
		assert.NoError(t, err)

		id, _ = result.LastInsertId()

		assert.NotEqual(t, int64(0), id)

	})

	t.Run("Find", func(t *testing.T) {
		query := `
			SELECT id, application_id, username, password_hash, created, updated 
			FROM users 
			WHERE id = ?
		`

		userFound := &model.User{}
		row := client.DB.QueryRow(query, int64(id))
		err := row.Scan(
			&userFound.ID,
			&userFound.ApplicationID,
			&userFound.Username,
			&userFound.PasswordHash,
			&userFound.Created,
			&userFound.Updated,
		)

		assert.NoError(t, err)
		assert.Equal(t, id, userFound.ID)
		assert.Equal(t, user.Username, userFound.Username)
	})

}
