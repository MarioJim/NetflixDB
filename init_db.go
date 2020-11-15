package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func initDB(ctx context.Context, client *mongo.Client) error {
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
	collection := client.Database("netflixdb").Collection("titles")
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
