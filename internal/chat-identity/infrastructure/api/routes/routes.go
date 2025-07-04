package routes

import (
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/infrastructure/api/controllers"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/middlewares"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/settings"
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	authMiddleWare := middlewares.NewJWTRSAMiddleware()
	config := settings.GetConfig()

	index := router.Group("/")
	{
		index.GET("/", func(c *gin.Context) {
			c.Data(200, "text/html; charset=utf-8", []byte(fmt.Sprintf("Welcome to %s V%s", config.AppName, config.AppVersion)))
		})

		index.GET("/ping", func(c *gin.Context) {
			c.Data(200, "text/html; charset=utf-8", []byte("pong"))
		})

		index.GET("/alive", func(c *gin.Context) {
			c.Data(200, "text/html; charset=utf-8", []byte("OK"))
		})

		index.GET("/ready", func(c *gin.Context) {
			c.Data(200, "text/html; charset=utf-8", []byte("OK"))
		})
	}

	v1 := router.Group("/api/v1")
	authController := controllers.AuthController{}
	authV1 := v1.Group("/auth")
	{
		authV1.POST("/login", authController.LoginV1)
		authV1.POST("/signup", authController.SignUpV1)
	}

	_ = authMiddleWare
}
