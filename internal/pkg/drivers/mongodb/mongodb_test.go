package mongodb

import (
	"context"
	"errors"
	"template-go-api/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test_connect(t *testing.T) {
	ogMongoConnect := mongoConnect
	testURI := "mongodb://localhost:27017"

	defer func() {
		mongoConnect = ogMongoConnect
	}()

	mongoConnectSuccess := func(_ context.Context, _ ...*options.ClientOptions) (*mongo.Client, error) {
		return &mongo.Client{}, nil
	}

	t.Run("Should return error if there is a problem conecting to database", func(t *testing.T) {
		mongoConnect = func(_ context.Context, _ ...*options.ClientOptions) (*mongo.Client, error) {
			return nil, errors.New("Boom mongo")
		}

		res, err := connect(testURI)
		assert.Nil(t, res)
		assert.Errorf(t, err, "Boom mongo")
	})

	t.Run("Should return the mongo client if everything goes well", func(t *testing.T) {
		mongoConnect = mongoConnectSuccess
		res, err := connect(testURI)
		assert.NotNil(t, res)
		assert.NoError(t, err)
	})
}

func Test_ping(t *testing.T) {

	t.Run("Should return error if ping fails", func(t *testing.T) {
		mockedClient := new(mocks.MongoClient)
		mockedClient.On("Ping", mock.Anything, mock.Anything).Return(errors.New("Boom ping"))

		err := ping(mockedClient)
		assert.Errorf(t, err, "Boom ping")
		mockedClient.AssertExpectations(t)
	})

	t.Run("Should return no error if ping is successful", func(t *testing.T) {
		mockedClient := new(mocks.MongoClient)
		mockedClient.On("Ping", mock.Anything, mock.Anything).Return(nil)

		err := ping(mockedClient)
		assert.NoError(t, err)
		mockedClient.AssertExpectations(t)
	})
}

func Test_getDatabase(t *testing.T) {
	t.Run("Should return the mongo database", func(t *testing.T) {
		mockedClient := new(mocks.MongoClient)
		mockedClient.On("Database", "test").Return(&mongo.Database{})

		db := getDatabase(mockedClient, "test")
		assert.Equal(t, &mongo.Database{}, db)
		mockedClient.AssertExpectations(t)
	})
}

func Test_ConnectDB(t *testing.T) {
	ogMongoClient := mongoClient
	ogMongoDB := mongoDatabse
	ogPing := _ping
	defer func() {
		mongoClient = ogMongoClient
		mongoDatabse = ogMongoDB
		_ping = ogPing
	}()

	pingSuccess := func(client MongoClient) error {
		return nil
	}

	t.Run("Should return ping result if mongoClient already exists", func(t *testing.T) {
		mongoClient = &mongo.Client{}
		_ping = func(client MongoClient) error {
			assert.Equal(t, mongoClient, client)
			return nil
		}

		err := ConnectDB("", "")
		assert.NoError(t, err)
	})

	t.Run("Should return error if connect fn returns error", func(t *testing.T) {
		mongoClient = nil
		_ping = pingSuccess

		_connect = func(uri string) (*mongo.Client, error) {
			return nil, errors.New("Boom connect")
		}

		err := ConnectDB("", "")
		assert.Errorf(t, err, "Boom connect")
	})

	t.Run("Should set corresponding client and database instances if everything goes well", func(t *testing.T) {
		testURI := "mongodb://localhost:27017"
		mongoClient = nil
		mongoDatabse = nil
		_ping = pingSuccess

		_connect = func(uri string) (*mongo.Client, error) {
			assert.Equal(t, testURI, uri)
			return &mongo.Client{}, nil
		}

		_getDatabase = func(_ MongoClient, dbName string) *mongo.Database {
			assert.Equal(t, "test-db", dbName)
			return &mongo.Database{}
		}

		err := ConnectDB(testURI, "test-db")
		assert.NoError(t, err)
		assert.NotNil(t, mongoClient)
		assert.NotNil(t, mongoDatabse)
	})
}

func Test_GetCollection(t *testing.T) {
	ogMongoDB := mongoDatabse

	defer func() {
		mongoDatabse = ogMongoDB
	}()

	t.Run("Should return error if mongo database instance has not been initialized", func(t *testing.T) {
		mongoDatabse = nil
		db, err := GetCollection("")
		assert.Nil(t, db)
		assert.Errorf(t, err, "Mongo connection to DB has not been established. Probably you need to call ConnectDB first")
	})
}
