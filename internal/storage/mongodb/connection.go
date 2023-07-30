package mongodb

import (
	"context"
	"myschool/pkg/helper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const EnvMongoURI string = "MONGODB_URI"
const DefaultMongoURI string = "mongodb://localhost:27017"

func Connect() (*mongo.Database, *mongo.Client, error) {

	mongoURI := helper.GetEnvString(EnvMongoURI, DefaultMongoURI)
	opt := options.Client().ApplyURI(mongoURI)
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
