package db

import (
	"bikagame-go/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DB struct {
	Client       *mongo.Client
	Database     *mongo.Database
	Users        *mongo.Collection
	Groups       *mongo.Collection
	Config       *mongo.Collection
	Transactions *mongo.Collection
	Orders       *mongo.Collection
}

func Connect(cfg config.Config) (*DB, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}
	database := client.Database(cfg.DBName)
	return &DB{
		Client:       client,
		Database:     database,
		Users:        database.Collection("users"),
		Groups:       database.Collection("groups"),
		Config:       database.Collection("config"),
		Transactions: database.Collection("transactions"),
		Orders:       database.Collection("orders"),
	}, nil
}
