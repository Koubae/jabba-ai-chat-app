package redis

import "github.com/Koubae/jabba-ai-chat-app/pkg/common/utils"

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	DB       int
	password string
	PoolSize int
}

func (c *DatabaseConfig) GetPass() string {
	return c.password
}

var databaseConfig *DatabaseConfig

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	host := utils.GetEnvString("REDIS_HOST", "localhost")
	port := utils.GetEnvInt("REDIS_PORT", 6379)
	db := utils.GetEnvInt("REDIS_DB", 0)
	password := utils.GetEnvString("REDIS_PASS", "")
	poolSize := utils.GetEnvInt("REDIS_MAX_CONNECTIONS", 50)

	config := &DatabaseConfig{
		Driver:   "redis",
		Host:     host,
		Port:     port,
		DB:       db,
		password: password,
		PoolSize: poolSize,
	}

	databaseConfig = config
	return config, nil
}

func GetDatabaseConfig() *DatabaseConfig {
	if databaseConfig == nil {
		panic("DatabaseConfig is not initialized, call LoadDatabaseConfig() first!")
	}
	return databaseConfig
}
