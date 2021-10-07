package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var mongoDatabse *mongo.Database

func connect(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func ping(client *mongo.Client) error {
	err := mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	return nil
}

func ConnectDB(uri string, dbName string) error {
	if mongoClient != nil {
		return ping(mongoClient)
	}

	client, err := connect(uri)
	if err != nil {
		return err
	}

	mongoClient = client
	mongoDatabse = mongoClient.Database(dbName)

	return ping(mongoClient)
}

func GetCollection(name string) (*mongo.Collection, error) {
	if mongoDatabse == nil {
		return nil, errors.New("Mongo connection to DB has not been established. Probably you need to call ConnectDB first")
	}

	return mongoDatabse.Collection(name), nil
}
