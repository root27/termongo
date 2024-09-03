package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type mongoClient struct {
	*mongo.Client
}

func connectToMongoDB() (*mongoClient, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(host))
	if err != nil {
		return nil, err
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &mongoClient{
		client}, nil
}

func (c *mongoClient) getCollections() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the list of collections
	collections, err := c.Database(dbName).ListCollectionNames(ctx, bson.D{})
	if err != nil {

		return nil, err
	}
	return collections, nil
}
