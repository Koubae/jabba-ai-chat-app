package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var client *Client

func NewClient() (*Client, error) {
	var err error
	var db *sql.DB
	var config *DatabaseConfig

	config, err = LoadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	db, err = sql.Open("mysql", config.Dns()+"?parseTime=true")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(int(config.maxOpenConnections))
	db.SetMaxIdleConns(int(config.maxIdleConnections))
	db.SetConnMaxLifetime(time.Duration(config.maxConnectionLifetime))
	db.SetConnMaxIdleTime(time.Duration(config.maxConnectionIdleTime))

	client = &Client{
		Config: config,
		DB:     db,
	}
	client.Ping()
	return client, nil

}

func GetClient() *Client {
	if client == nil {
		panic("MySQL Client is not initialized!")
	}
	return client
}

type Client struct {
	Config *DatabaseConfig
	DB     *sql.DB
}

func (c *Client) String() string {
	return fmt.Sprintf("Client{config: %+v}", c.Config)
}

func (c *Client) Shutdown() {
	err := c.DB.Close()
	if err != nil {
		log.Printf("failed to close MySQL: %v\n", err.Error())
		return
	}
}

func (c *Client) Ping() {
	err := c.DB.Ping()
	if err != nil {
		log.Fatalf("failed to ping MySQL: %v\n", err.Error())
	}

}
