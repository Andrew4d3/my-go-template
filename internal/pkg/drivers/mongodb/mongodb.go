package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client
var mongoDatabse *mongo.Database

// MongoClient defines a mongo client interface
type MongoClient interface {
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
}

var mongoConnect = mongo.Connect

func connect(uri string) (*mongo.Client, error) {
	client, err := mongoConnect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func ping(client MongoClient) error {
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	return nil
}

func getDatabase(client MongoClient, dbName string) *mongo.Database {
	return client.Database(dbName)
}

var _connect = connect
var _ping = ping
var _getDatabase = getDatabase

// ConnectDB establishes a mongo db connection
func ConnectDB(uri string, dbName string) error {
	if mongoClient != nil {
		return _ping(mongoClient)
	}

	client, err := _connect(uri)
	if err != nil {
		return err
	}

	mongoClient = client
	mongoDatabse = _getDatabase(client, dbName)

	return _ping(mongoClient)
}

// GetCollection gets a collection from the current connection
func GetCollection(name string) (*mongo.Collection, error) {
	if mongoDatabse == nil {
		return nil, errors.New("Mongo connection to DB has not been established. Probably you need to call ConnectDB first")
	}

	// Not possible to unit test. Consider wrapping this around interface
	return mongoDatabse.Collection(name), nil
}
