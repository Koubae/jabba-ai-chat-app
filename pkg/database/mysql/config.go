package mysql

import (
	"fmt"
	"github.com/Koubae/jabba-ai-chat-app/pkg/common/utils"
)

type DatabaseConfig struct {
	Driver                string
	DBName                string
	host                  string
	port                  uint16
	user                  string
	password              string
	maxOpenConnections    uint
	maxIdleConnections    uint
	maxConnectionLifetime uint
	maxConnectionIdleTime uint
}

var databaseConfig *DatabaseConfig

func (c *DatabaseConfig) Dns() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", c.user, c.password, c.host, c.port, c.DBName)
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	name := utils.GetEnvString("MYSQL_DB", "")
	host := utils.GetEnvString("MYSQL_HOST", "localhost")
	port := utils.GetEnvInt("MYSQL_PORT", 3306)
	user := utils.GetEnvString("MYSQL_USER", "root")
	password := utils.GetEnvString("MYSQL_PASS", "admin")

	maxOpenConnections := utils.GetEnvInt("MYSQL_MAX_OPEN_CONNECTIONS", 50)
	maxIdleConnections := utils.GetEnvInt("MYSQL_MAX_IDLE_CONNECTIONS", 25)
	maxConnectionLifetime := utils.GetEnvInt("MYSQL_MAX_CONNECTION_LIFETIME", 300)
	maxConnectionIdleTime := utils.GetEnvInt("MYSQL_MAX_CONNECTION_IDLE_TIME", 600)

	config := &DatabaseConfig{
		Driver:                "mysql",
		DBName:                name,
		host:                  host,
		port:                  uint16(port),
		user:                  user,
		password:              password,
		maxOpenConnections:    uint(maxOpenConnections),
		maxIdleConnections:    uint(maxIdleConnections),
		maxConnectionLifetime: uint(maxConnectionLifetime),
		maxConnectionIdleTime: uint(maxConnectionIdleTime),
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
