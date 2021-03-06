package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	database = "btcPrice"
	// CollectionPrice a collection name for CRU pic
	CollectionPrice = "prices"
	// CollectionUser a collection name for CR user
	CollectionUser = "users"
	// CollectionLatistedID a collection name for foreign key(?) for prices
	CollectionLatistedID = "latestID"
)

var mongoClient *mongo.Client

// SetMongoClint initialize mongoClient
func SetMongoClint(c *mongo.Client) {
	// can be more explicitly
	if mongoClient == nil {
		mongoClient = c // <--- NOT THREAD SAFE
	}
}

// GetCollection to get the connection for mongodb collection
func GetCollection(collectionName string) *mongo.Collection {
	return mongoClient.Database(database).Collection(collectionName)
}

// Create one obj into specify collection
func Create(collectionName string, item interface{}) (*mongo.InsertOneResult, error) {
	collection := mongoClient.Database(database).Collection(collectionName)
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	return collection.InsertOne(ctx, item)
}

// Delete one obj from specify collection
func Delete(collectionName string, _id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": bson.M{"$eq": _id}}
	collection := mongoClient.Database(database).Collection(collectionName)
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	return collection.DeleteOne(ctx, filter)
}

// Update one obj from specify collection
func Update(collectionName string, _id primitive.ObjectID, item interface{}) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": bson.M{"$eq": _id}}
	update := bson.M{"$set": item}
	collection := mongoClient.Database(database).Collection(collectionName)
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	return collection.UpdateMany(ctx, filter, update)
}

// FindOne find one obj from specify collection
func FindOne(collectionName string, filter interface{}) (r *mongo.SingleResult) {
	collection := mongoClient.Database(database).Collection(collectionName)
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	return collection.FindOne(ctx, filter)
}
