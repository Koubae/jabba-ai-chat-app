package mongodb

import (
	"errors"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/utils"
)

type DatabaseConfig struct {
	Driver string
	Uri    string
	DBName string
}

var databaseConfig *DatabaseConfig

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	databaseUri := utils.GetEnvString("MONGODB_URI", "")
	name := utils.GetEnvString("MONGODB_DB", "")
	if databaseUri == "" || name == "" {
		return nil, errors.New("database configuration missing")
	}

	config := &DatabaseConfig{
		Driver: "mongodb",
		Uri:    databaseUri,
		DBName: name,
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
