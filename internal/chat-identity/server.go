package chat_identity

import (
	"context"
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/container"
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity/infrastructure/api/routes"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/settings"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	err := godotenv.Load(".env.chat-identity")
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
