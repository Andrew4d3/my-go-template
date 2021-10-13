package repository

// NOTE: Not possible to unit test. Mongo driver uses structs instead of interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository defines a mongo repository
type MongoRepository struct {
	collection Collection
}

// Collection provides a mongo collection interface
type Collection interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
}

// New creates a new mongo repository instance
func New(collection Collection) *MongoRepository {
	return &MongoRepository{collection: collection}
}

// FindOne gets a unique collection element based on a query
func (m MongoRepository) FindOne(filter bson.M, result interface{}) error {
	err := m.collection.FindOne(context.TODO(), filter).Decode(result)

	if err != nil && err == mongo.ErrNoDocuments {
		return nil
	}

	return err
}
