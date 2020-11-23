package main

import (
	"bufio"
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Query : Entry point for querying a document from the database
func Query(ctx context.Context, client *mongo.Client, scanner *bufio.Scanner, rdb *redis.Client) error {
	fmt.Println("Query")
	fmt.Println()
	fmt.Println("	1) Query for Movies")
	fmt.Println("	2) Query for Tv Shows")
	fmt.Println("	3) Query for actors")
	actionStr := ScanStringWithPrompt("Select an action: ", scanner)
	action, err := strconv.Atoi(actionStr)
	if err != nil {
		action = 0
	}
	switch action {
	case 1:
		title := ScanStringWithPrompt("Write the title of the movie you want to query: ", scanner)
		key := "Movie" + title
		val, err := rdb.Get(ctx, key).Result()
		if err == nil {
			fmt.Println("Found in cache")
			fmt.Println(val)
			return nil
		}
		val = ""
		filter := bson.D{primitive.E{Key: "title", Value: title},
			primitive.E{Key: "type", Value: "Movie"}}
		opts := options.FindOne().SetProjection(bson.M{"_id": 0, "director": 1, "cast": 1, "country": 1, "release_year": 1})
		var searchResult bson.M
		err = GetTitlesColl(client).FindOne(ctx, filter, opts).Decode(&searchResult)
		if err != nil && err == mongo.ErrNoDocuments {
			fmt.Println("A movie with that title wasn't found")
			return nil
		} else if err != nil {
			return err
		}
		for key, value := range searchResult {
			val += fmt.Sprintf("%s: %s\n", key, ValueToString(value))

		}
		fmt.Println(val)

		err = rdb.Set(ctx, key, val, 0).Err()
		return err

	case 2:
		title := ScanStringWithPrompt("Write the title of the tv show you want to query: ", scanner)
		key := "TV Show" + title
		val, err := rdb.Get(ctx, key).Result()
		if err == nil {
			fmt.Println("Found in cache")
			fmt.Println(val)
			return nil
		}
		val = ""
		filter := bson.D{primitive.E{Key: "title", Value: title},
			primitive.E{Key: "type", Value: "TV Show"}}
		opts := options.FindOne().SetProjection(bson.M{"_id": 0, "director": 1, "cast": 1, "country": 1, "release_year": 1})
		var searchResult bson.M
		err = GetTitlesColl(client).FindOne(ctx, filter, opts).Decode(&searchResult)
		if err != nil && err == mongo.ErrNoDocuments {
			fmt.Println("A tv show with that title wasn't found")
			return nil
		} else if err != nil {
			return err
		}
		for key, value := range searchResult {
			val += fmt.Sprintf("%s: %s\n", key, ValueToString(value))
		}
		fmt.Println(val)

		err = rdb.Set(ctx, key, val, 0).Err()
		return err

	case 3:
		actor := ScanStringWithPrompt("Write the name of the actor you want to query: ", scanner)
		key := "TV Show" + actor
		val, err := rdb.Get(ctx, key).Result()
		if err == nil {
			fmt.Println("Found in cache")
			fmt.Println(val)
			return nil
		}
		val = ""
		unwind := bson.D{primitive.E{Key: "$unwind", Value: "$cast"}}
		matchMovies := bson.D{primitive.E{Key: "$match", Value: bson.D{
			primitive.E{Key: "type", Value: "Movie"},
		}}}
		matchTV := bson.D{primitive.E{Key: "$match", Value: bson.D{
			primitive.E{Key: "type", Value: "TV Show"},
		}}}
		matchActor := bson.D{primitive.E{Key: "$match", Value: bson.D{
			primitive.E{Key: "cast", Value: actor},
		}}}
		project := bson.D{primitive.E{Key: "$project", Value: bson.D{
			primitive.E{Key: "_id", Value: 0},
			primitive.E{Key: "title", Value: 1},
		}}}
		cursor, err := GetTitlesColl(client).Aggregate(ctx, mongo.Pipeline{matchMovies, unwind, matchActor, project})
		if err != nil {
			return err
		}
		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			return err
		}
		val += fmt.Sprintln("Movies:")
		for _, result := range results {
			val += fmt.Sprintf(" - %s\n", result["title"])
		}
		cursor, err = GetTitlesColl(client).Aggregate(ctx, mongo.Pipeline{matchTV, unwind, matchActor, project})
		if err != nil {
			return err
		}

		if err = cursor.All(ctx, &results); err != nil {
			return err
		}
		val += fmt.Sprintln("TV Shows:")
		for _, result := range results {
			val += fmt.Sprintf(" - %s\n", result["title"])
		}
		fmt.Println(val)
		err = rdb.Set(ctx, key, val, 0).Err()
		return err
	default:
		fmt.Println("Action couldn't be identified")
	}
	return nil
}
