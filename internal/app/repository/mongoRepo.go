package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func New(collection *mongo.Collection) *MongoRepository {
	return &MongoRepository{collection: collection}
}

func (m MongoRepository) FindOne(filter bson.M, result interface{}) error {
	err := m.collection.FindOne(context.TODO(), filter).Decode(result)

	if err != nil && err == mongo.ErrNoDocuments {
		return nil
	}

	return err
}
