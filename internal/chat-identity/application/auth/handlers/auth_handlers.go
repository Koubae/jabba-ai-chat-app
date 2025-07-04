package handlers

import (
	"crypto/rsa"
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/domain/user/model"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/settings"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/utils"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

type LoginRequest struct {
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
	ApplicationID string `json:"application_id" binding:"required"`
}

func (r *LoginRequest) Validate() error {
	r.Username = strings.TrimSpace(r.Username)
	r.Password = strings.TrimSpace(r.Password)
	r.ApplicationID = strings.TrimSpace(r.ApplicationID)

	if r.Username == "" {
		return errors.New("username is required")
	} else if r.Password == "" {
		return errors.New("password is required")
	} else if r.ApplicationID == "" {
		return errors.New("application_id is required")
	}

	return nil
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Expires     int64  `json:"expires"`
}

type LoginHandler struct {
	Command  LoginRequest
	Response LoginResponse
}

func (h *LoginHandler) Handle() error {
	expire := time.Now().Add(settings.AuthTokenExpirationTime).Unix()

	// TODO Mocking find userID in db | Add real database here!
	userID := uint(1)

	token, err := generateJWTWithRSA(userID, h.Command.Username, h.Command.ApplicationID, expire)
	if err != nil {
		return err
	}

	h.Response = LoginResponse{AccessToken: token, Expires: expire}
	return nil
}

func generateJWTWithRSA(userID uint, userName string, ApplicationID string, expire int64) (string, error) {
	privateKey := loadAndGetPrivateKey()
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": expire,
		"iss": "jabba-ai-chat",

		"role":           "user",
		"user_name":      userName,
		"application_id": ApplicationID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)

}

var privateKey *rsa.PrivateKey

func loadAndGetPrivateKey() *rsa.PrivateKey {
	var err error
	if privateKey == nil {
		privateKey, err = utils.GetPrivateKey()
		if err != nil {
			panic(err.Error())
		}
	}
	return privateKey
}

type SignUpRequest struct {
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
	ApplicationID string `json:"application_id" binding:"required"`
}

func (r *SignUpRequest) Validate() error {
	// Normalize Data
	r.Username = strings.TrimSpace(r.Username)
	r.Password = strings.TrimSpace(r.Password)
	r.ApplicationID = strings.TrimSpace(r.ApplicationID)

	if r.Username == "" {
		return errors.New("username is required")
	} else if r.Password == "" {
		return errors.New("password is required")
	} else if r.ApplicationID == "" {
		return errors.New("application_id is required")
	}

	return nil
}

type SignUpResponse struct {
	model.User
}

type SignUpHandler struct {
	Command  SignUpRequest
	Response SignUpResponse
}

func (h *SignUpHandler) Handle() error {
	// TODO: check that user DOES NMOT exists in db
	ID := uint(1) // TODO: Create from db for real
	user := model.User{
		ID:            int64(ID),
		Username:      h.Command.Username,
		ApplicationID: h.Command.ApplicationID,
	}

	h.Response = SignUpResponse{
		User: user,
	}

	return nil
}
