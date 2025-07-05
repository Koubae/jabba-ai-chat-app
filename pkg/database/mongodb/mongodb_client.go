package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var dbClient *Client

func NewClient() (*Client, error) {
	var err error
	var client *mongo.Client
	var config *DatabaseConfig

	config, err = LoadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opt := options.Client().
		ApplyURI(config.Uri).
		SetMinPoolSize(5).
		SetMaxPoolSize(100).
		SetMaxConnIdleTime(10 * time.Minute)

	client, err = mongo.Connect(ctx, opt)
	if err != nil {
		return nil, err
	}

	db := client.Database(config.DBName)
	dbClient = &Client{
		config: config,
		client: client,
		db:     db,
	}
	dbClient.Ping(ctx)
	return dbClient, nil
}

func GetClient() *Client {
	if dbClient == nil {
		panic("MongoDB Client is not initialized!")
	}
	return dbClient
}

type Client struct {
	config *DatabaseConfig
	client *mongo.Client
	db     *mongo.Database
}

func (c *Client) String() string {
	return fmt.Sprintf("Client{config: %v}", c.config.DBName)
}

func (c *Client) Shutdown(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("error disconnecting from MongoDB: %w\n", err)
	}
	return nil
}

func (c *Client) Ping(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := c.client.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping MongoDB: %v\n", err.Error())
	}
}

func (c *Client) Collection(name string) *mongo.Collection {
	return c.db.Collection(name)
}

func (c *Client) ListDatabases(ctx context.Context) ([]string, error) {
	databases, err := c.client.ListDatabaseNames(ctx, bson.D{{"empty", false}})
	if err != nil {
		return nil, err
	}
	return databases, nil
}

func (c *Client) CreateUniqueIndex(collection *mongo.Collection, ctx context.Context, field string) error {
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{field: 1},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func (c *Client) CreateCompoundUniqueIndex(collection *mongo.Collection, ctx context.Context, fields []string) error {
	keys := bson.D{}
	for _, field := range fields {
		keys = append(keys, bson.E{Key: field, Value: 1})
	}

	indexModel := mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

// CreateIndex Usage examples:
// Single field ascending
// err = client.CreateIndex(collection, context.Background(), []string{"user_id"}, []int{1})
// Multiple fields with mixed order
// err = client.CreateIndex(collection, context.Background(), []string{"user_id", "created_at"}, []int{1, -1})
func (c *Client) CreateIndex(collection *mongo.Collection, ctx context.Context, fields []string, orders []int) error {
	if len(fields) != len(orders) {
		return fmt.Errorf("fields and orders must have the same length")
	}

	keys := bson.D{}
	for i, field := range fields {
		keys = append(keys, bson.E{Key: field, Value: orders[i]})
	}

	indexModel := mongo.IndexModel{
		Keys: keys,
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}
