package chat_session

import (
	"context"
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/container"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-session/infrastructure/api/routes"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/settings"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	err := godotenv.Load(".env.chat-session")
	if err != nil {
		panic(err.Error())
	}
	config := settings.NewConfig()
	switch config.Environment {
	case settings.EnvTesting:
		gin.SetMode(gin.TestMode)
	case settings.EnvDev, settings.EnvStaging:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	container.CreateDIContainer()

}

func RunServer() {
	config := settings.GetConfig()

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	err := router.SetTrustedProxies(config.TrustedProxies)
	if err != nil {
		panic(err.Error())
	}
	routes.InitRoutes(router)

	srv := &http.Server{
		Addr:    config.GetAddr(),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error while shutting down server, error: %s\n", err)
		}
	}()
	log.Printf("Service %s V%s | Server started on port %s\n", config.AppName, config.AppVersion, config.GetAddr())

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify the user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server Shutdown, cleaning up resources")
	container.ShutDown()

	log.Println("Server exiting")
}
