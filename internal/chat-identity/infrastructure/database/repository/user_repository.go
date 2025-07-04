package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/domain/user/model"
	domainrepository "github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/domain/user/repository"
	"github.com/Koubae/jabba-ai-chat-app/pkg/database/mysql"
	mysqladapter "github.com/go-sql-driver/mysql"
)

type UserRepository struct {
	db *mysql.Client
}

func NewUserRepository(db *mysql.Client) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (application_id, username, password_hash) 
		VALUES (?, ?, ?)
	`
	result, err := r.db.DB.Exec(query, user.ApplicationID, user.Username, user.PasswordHash)
	if err != nil {
		var mysqlErr *mysqladapter.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 { // MySQL error 1062 = ER_DUP_ENTRY
			return domainrepository.ErrUserAlreadyExists
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil

}

func (r *UserRepository) GetByID(id int64) (*model.User, error) {
	query := `
		SELECT id, application_id, username, password_hash, created, updated 
		FROM users 
		WHERE id = ?
	`

	user := &model.User{}
	row := r.db.DB.QueryRow(query, id)
	err := row.Scan(
		&user.ID,
		&user.ApplicationID,
		&user.Username,
		&user.PasswordHash,
		&user.Created,
		&user.Updated,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainrepository.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil

}

func (r *UserRepository) GetByUsername(applicationID string, username string) (*model.User, error) {
	query := `
		SELECT id, application_id, username, password_hash, created, updated 
		FROM users 
		WHERE application_id = ? AND username = ?
	`

	user := &model.User{}
	row := r.db.DB.QueryRow(query, applicationID, username)
	err := row.Scan(
		&user.ID,
		&user.ApplicationID,
		&user.Username,
		&user.PasswordHash,
		&user.Created,
		&user.Updated,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainrepository.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil

}
