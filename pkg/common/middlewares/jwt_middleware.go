package middlewares

import (
	"crypto/rsa"
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type JWTSecret interface {
	[]byte | *rsa.PublicKey | *rsa.PrivateKey
}

func jwtMiddleware[S JWTSecret](c *gin.Context, method jwt.SigningMethod, secret S) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method != method {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized -- " + err.Error()})
		return
	} else if !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	applicationId := claims["application_id"].(string)
	userId := int64(claims["sub"].(float64))
	issuer := claims["iss"].(string)
	role := claims["role"].(string)
	userName := claims["user_name"].(string)
	accessToken := &auth.AccessToken{
		ApplicationId: claims["application_id"].(string),
		UserId:        int64(claims["sub"].(float64)),
		Username:      claims["user_name"].(string),
		Issuer:        claims["iss"].(string),
		Role:          claims["role"].(string),
		AccessToken:   tokenString,
	}

	c.Set("application_id", applicationId)
	c.Set("user_id", userId)
	c.Set("issuer", issuer)
	c.Set("role", role)
	c.Set("user_name", userName)
	c.Set("access_token", accessToken)

	c.Next()
}
