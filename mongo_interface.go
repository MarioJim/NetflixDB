package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoURI string = "mongodb://localhost:27017"

// ConnectMongoDB : Setups a connection to the MongoDB server
func ConnectMongoDB() (context.Context, *mongo.Client, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	return ctx, client, cancel, err
}

// GetTitlesColl : Returns a reference to the collection of titles
func GetTitlesColl(client *mongo.Client) *mongo.Collection {
	return client.Database("netflixdb").Collection("titles")
}

// InitMongoDB : Initializes the database with documents from a JSON file
func InitMongoDB(ctx context.Context, client *mongo.Client) error {
	fmt.Println("Checking if the MongoDB collection is initialized...")
	result, err := client.ListDatabaseNames(ctx, bson.D{primitive.E{
		Key:   "name",
		Value: "netflixdb",
	}})
	if err != nil {
		return err
	}
	if len(result) == 1 {
		fmt.Println("It is!")
		return nil
	}
	fmt.Println("Collection not found, reading documents from 'dataset/netflix_titles.json'...")
	byteValues, err := ioutil.ReadFile("dataset/netflix_titles.json")
	if err != nil {
		return err
	}
	var docs []NetflixTitle
	err = json.Unmarshal(byteValues, &docs)
	if err != nil {
		return err
	}
	collection := GetTitlesColl(client)
	for i := range docs {
		doc := docs[i]
		_, err := collection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}
	}
	fmt.Println("Documents inserted successfully!")
	return nil
}
