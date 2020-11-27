package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Add : Entry point for adding a new document to the database
func Add(ctx context.Context, client *mongo.Client, scanner *bufio.Scanner, rdb *redis.Client) error {
	dummyBson := bson.M{
		"cast":         primitive.A{},
		"country":      primitive.A{},
		"date_added":   "",
		"description":  "",
		"director":     primitive.A{},
		"duration":     "",
		"listed_in":    primitive.A{},
		"rating":       "",
		"release_year": 0.0,
		"title":        "",
		"type":         "",
	}

	var randomID int
	for {
		randomID = rand.Intn(100000000)
		filter := bson.D{primitive.E{Key: "show_id", Value: randomID}}
		var checkResult bson.M
		err := GetTitlesColl(client).FindOne(ctx, filter).Decode(&checkResult)
		if err != nil && err == mongo.ErrNoDocuments {
			break
		}
	}

	newFields := bson.D{primitive.E{Key: "show_id", Value: randomID}}
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
	fmt.Printf("Correctly added document %d!\n", randomID)

	return rdb.FlushDB(ctx).Err()
}
