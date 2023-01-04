package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//TODO get from .env
// Database settings (insert your own database name and connection URI)
const dbName = "users"
const mongoURI = "mongodb://mongouser:mongoPass@mongodb:27017/users?authSource=admin"

// MongoInstance contains the Mongo client and database objects
type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var Instance MongoInstance

func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {

		return err
	}

	Instance = MongoInstance{
		Client: client,
		Db:     db,
	}

	return nil
}
