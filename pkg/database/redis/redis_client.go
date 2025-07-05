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

func GetClient() *Client {
	if client == nil {
		panic("Redis Client is not initialized!")
	}
	return client
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

// AllKeys
// /!\/!\/!\/!\/!\ Use this ONLY for debugging!
func (c *Client) AllKeys(ctx context.Context) ([]string, error) {
	keys, err := c.DB.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}
	return keys, nil

}

// LogAllValues
// /!\/!\/!\/!\/!\ Use this ONLY for debugging!
func (c *Client) LogAllValues(ctx context.Context) error {
	keys, err := c.AllKeys(ctx)
	if err != nil {
		return err
	}

	for _, key := range keys {
		value, err := c.DB.Get(ctx, key).Result()
		if err != nil {
			fmt.Printf("Error getting key %s: %v\n", key, err)
			return err
		}
		fmt.Printf("-Key: %s\n-Value:\n%s\n------\n", key, value)
	}
	return nil
}
