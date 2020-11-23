package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Add : Entry point for adding a new document to the database
func Add(ctx context.Context, client *mongo.Client, scanner *bufio.Scanner) error {
	dummyBson := bson.M{
		"cast":primitive.A{},
		"country":primitive.A{}, 
		"date_added": "", 
		"description": "When a prison ship crash unleashes hundreds of Decepticons on Earth, Bumblebee leads a new Autobot force to protect humankind.",
		"director": primitive.A{},
		"duration": "", 
		"listed_in": primitive.A{}, 
		"rating": "", 
		"release_year": 0.0,
		"title": "",
		"type": "",
	}

	var randomId int
	checkIfRepeated := true
	for checkIfRepeated {
		randomId = rand.Intn(99999999)
		filter := bson.D{primitive.E{Key: "show_id", Value: randomId}}
		var checkResult bson.M
		err := GetTitlesColl(client).FindOne(ctx, filter).Decode(&checkResult)
		if err != nil && err == mongo.ErrNoDocuments {
			checkIfRepeated = false
		}
	}

	newFields := bson.D{primitive.E{Key:"show_id", Value: randomId}}
	for key, value := range dummyBson {
		inpMessage := fmt.Sprintf("Insert value for '%s': ", key)
		switch value.(type) {
		case string:
			newFields = append(newFields, primitive.E{
				Key:   key,
				Value: ScanStringWithPrompt(inpMessage, scanner),
			})
		case float64:
			newFields = append(newFields, primitive.E{
				Key:   key,
				Value: ScanNumberWithPrompt(inpMessage, scanner),
			})
		case primitive.A:
			fmt.Printf("Insert values for '%s':\n", key)
			newArray := append(primitive.A{}, ScanStringArray(scanner)...)
			newFields = append(newFields, primitive.E{
				Key:   key,
				Value: newArray,
			})
		}
	}

	if _, err := GetTitlesColl(client).InsertOne(ctx, newFields); err != nil {
		return err
	}
	fmt.Printf("Correctly updated document %d!\n", randomId)
	return nil
}
