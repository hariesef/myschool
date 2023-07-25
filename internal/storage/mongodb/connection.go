package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect() (*mongo.Database, *mongo.Client, error) {

	opt := options.Client().ApplyURI("mongodb://localhost:27017")
	localMongoClient, err := mongo.Connect(context.Background(), opt)
	if err != nil {
		return nil, nil, err
	}

	err = localMongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, nil, err
	}

	db := localMongoClient.Database("main")
	return db, localMongoClient, nil

}
