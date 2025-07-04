package settings

import (
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/utils"
	"os"
	"slices"
	"strconv"
)

type Config struct {
	port           uint16
	AppName        string
	AppVersion     string
	Environment    string
	TrustedProxies []string
}

func (c Config) GetAddr() string {
	return fmt.Sprintf(":%d", c.port)
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		panic("Config is not initialized, call NewConfig() first!")
	}
	return config
}
func NewConfig() *Config {
	port := utils.GetEnvInt("APP_PORT", 8001)

	errTemp := os.Setenv("PORT", strconv.Itoa(port)) // For gin-gonic
	if errTemp != nil {
		panic(errTemp.Error())
	}

	appName := utils.GetEnvString("APP_NAME", "unknown")
	appVersion := utils.GetEnvString("APP_VERSION", "unknown")

	environment := utils.GetEnvString("APP_ENVIRONMENT", "development")
	if !slices.Contains(Environments[:], environment) {
		panic(fmt.Sprintf("Invalid environment: %s, supported envs are %v", environment, Environments))
	}
	trustedProxies := utils.GetEnvStringSlice("APP_NETWORKING_PROXIES", []string{})

	config = &Config{
		port:           uint16(port),
		AppName:        appName,
		AppVersion:     appVersion,
		Environment:    environment,
		TrustedProxies: trustedProxies,
	}
	return config
}
