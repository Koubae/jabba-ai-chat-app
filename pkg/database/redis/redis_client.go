package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

var client *Client

func NewClient() (*Client, error) {
	var err error
	var config *DatabaseConfig

	config, err = LoadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.GetPass(), // no password set
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	client = &Client{
		Config: config,
		DB:     db,
	}
	client.Ping(ctx)
	return client, nil

}

type Client struct {
	Config *DatabaseConfig
	DB     *redis.Client
}

func (c *Client) String() string {
	return fmt.Sprintf("Client{config: %+v}", c.Config)
}

func (c *Client) Shutdown() {
	if err := c.DB.Close(); err != nil {
		log.Printf("failed to shutdown Redis: %v\n", err.Error())
		return
	}
}

func (c *Client) Ping(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := c.DB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to ping Redis: %v\n", err.Error())
	}
}
