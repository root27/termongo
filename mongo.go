package main

import (
	"context"
	"encoding/json"
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

func (c *mongoClient) numberOfDocs(collection string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Get the number of documents in the collection
	count, err := c.Database(dbName).Collection(collection).CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *mongoClient) getCollSize(collection string) (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Get the size of the collection
	stats := c.Database(dbName).RunCommand(ctx, bson.D{{"collStats", collection}})

	var result bson.M

	err := stats.Decode(&result)

	if err != nil {
		return bson.M{}, err
	}

	return result, nil
}

func (c *mongoClient) findAll(collection string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Find all documents in the collection
	cursor, err := c.Database(dbName).Collection(collection).Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	var jsonData []byte

	for _, result := range results {
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return nil, err
		}

		jsonData = append(jsonData, data...)

	}
	return jsonData, nil
}

func (c *mongoClient) findOne(collection string, filter string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var mongoFilter bson.M

	json.Unmarshal([]byte(filter), &mongoFilter)

	// Find one document in the collection
	var result bson.M
	err := c.Database(dbName).Collection(collection).FindOne(ctx, mongoFilter).Decode(&result)
	if err != nil {
		return nil, err
	}

	data, err := json.MarshalIndent(result, "", "  ")

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *mongoClient) findOneAndUpdate(collection string, filter string, update string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var mongoFilter bson.M
	json.Unmarshal([]byte(filter), &mongoFilter)
	var mongoUpdate bson.M
	json.Unmarshal([]byte(update), &mongoUpdate)
	// Find one document in the collection and update it

	var updateResult bson.M
	err := c.Database(dbName).Collection(collection).FindOneAndUpdate(ctx, mongoFilter, mongoUpdate,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updateResult)
	if err != nil {
		return nil, err
	}

	data, err := json.MarshalIndent(updateResult, "", "  ")

	if err != nil {
		return nil, err
	}

	return data, nil

}

func (c *mongoClient) insertOne(collection string, document string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var mongoDocument bson.M
	json.Unmarshal([]byte(document), &mongoDocument)
	// Insert one document into the collection
	_, err := c.Database(dbName).Collection(collection).InsertOne(ctx, mongoDocument)
	if err != nil {
		return nil, err
	}
	return []byte("Document inserted successfully"), nil
}

func (c *mongoClient) deleteOne(collection string, filter string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var mongoFilter bson.M
	json.Unmarshal([]byte(filter), &mongoFilter)
	// Delete one document from the collection
	_, err := c.Database(dbName).Collection(collection).DeleteOne(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	return []byte("Document deleted successfully"), nil
}
