package middlewares

import (
	"crypto/rsa"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func NewJWTRSAMiddleware() gin.HandlerFunc {
	publicKey := utils.GetPublicKeyOrPanic()
	return func(c *gin.Context) {
		jwtMiddleware[*rsa.PublicKey](c, jwt.SigningMethodRS256, publicKey)
	}
}
