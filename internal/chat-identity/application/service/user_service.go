package service

import (
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/domain/user/model"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/domain/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) CreateUser(applicationID string, username string, password string) (*model.User, error) {
	var user *model.User

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user = &model.User{
		ApplicationID: applicationID,
		Username:      username,
		PasswordHash:  string(passwordHash),
	}
	err = s.repository.Create(user)
	if err != nil {
		return nil, err
	}

	user, err = s.repository.GetByUsername(applicationID, username)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (s *UserService) GetUser(applicationID string, username string) (*model.User, error) {
	user, err := s.repository.GetByUsername(applicationID, username)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (s *UserService) VerifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
