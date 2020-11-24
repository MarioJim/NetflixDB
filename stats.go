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
	//"go.mongodb.org/mongo-driver/mongo/options"
)

// Statistics : Entry point for getting statistics from the database's documents
func Statistics (ctx context.Context, client *mongo.Client, scanner *bufio.Scanner, rdb *redis.Client) error {
	
	PrintMenuInConsole()
	
	actionStr := ScanStringWithPrompt("Select an action: ", scanner)
	action, err := strconv.Atoi(actionStr)
	if err != nil {
		action = 0
	}
	
	switch action {
	//Count number of movies and TV shows in the database
	case 1:
		key := "count_movies_shows"
		val, err := rdb.Get(ctx, key).Result()
		if err == nil {
			fmt.Println("Found in cache")
			fmt.Println(val)
			return nil
		}
		val = ""
		matchMovies := bson.D{primitive.E{Key: "$match", Value: bson.D{
			primitive.E{Key: "type", Value: "Movie"},
		}}}
		matchTV := bson.D{primitive.E{Key: "$match", Value: bson.D{
			primitive.E{Key: "type", Value: "TV Show"},
		}}}
		groupMovies := bson.D{primitive.E{Key: "$group", Value: bson.D{
			primitive.E{Key: "_id", Value: "Movies"},
			primitive.E{Key: "count", Value: bson.D{
				primitive.E{Key: "$sum", Value: 1}}},
		}}}
		groupTVShows := bson.D{primitive.E{Key: "$group", Value: bson.D{
			primitive.E{Key: "_id", Value: "TV Shows"},
			primitive.E{Key: "count", Value: bson.D{
				primitive.E{Key: "$sum", Value: 1}}},
		}}}
		
		cursor, err := GetTitlesColl(client).Aggregate(ctx, mongo.Pipeline{matchMovies, groupMovies})
		
		if err != nil {
			return err
		}
		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			return err
		}
		val += fmt.Sprintln("Movies:")
		for _, result := range results {
			val += fmt.Sprintf(" - %d\n", result["count"])
		}
		cursor, err = GetTitlesColl(client).Aggregate(ctx, mongo.Pipeline{matchTV, groupTVShows})
		if err != nil {
			return err
		}

		if err = cursor.All(ctx, &results); err != nil {
			return err
		}
		val += fmt.Sprintln("TV Shows:")
		for _, result := range results {
			val += fmt.Sprintf(" - %d\n", result["count"])
		}
		
		fmt.Println(val)
		err = rdb.Set(ctx, key, val, 0).Err()
		return err
		
	//Total number of movies for a given country
	case 2:
		country := ScanStringWithPrompt("Write a country to count the total number of movies: ", scanner)
		key := "total_movies_" + country
		val, err := rdb.Get(ctx, key).Result()
		if err == nil {
			fmt.Println("Found in cache")
			fmt.Println(val)
			return nil
		}
		val = ""
		matchMovies := bson.D{primitive.E{Key: "$match", Value: bson.D{
			primitive.E{Key: "type", Value: "Movie"},
		}}}
		unwindCountries := bson.D{primitive.E{Key: "$unwind", Value: "$country"}}
		matchCountries := bson.D{primitive.E{Key: "$match", Value: bson.D{
			primitive.E{Key: "$country", Value: country},
		}}}
		group := bson.D{primitive.E{Key: "$group", Value: bson.D{
			primitive.E{Key: "_id", Value: "total_movies"},
			primitive.E{Key: "total", Value: bson.D{
				primitive.E{Key: "$sum", Value: 1}}},
		}}}
		
		cursor, err := GetTitlesColl(client).Aggregate(ctx, mongo.Pipeline{matchMovies, unwindCountries, matchCountries, group})
		
		if err != nil {
			return err
		}
		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			return err
		}
		val += fmt.Sprintf("Total number of movies in %s:\n", country)
		for _, result := range results {
			val += fmt.Sprintf(" - %d\n", result["total"])
		}
		
		fmt.Println(val)
		err = rdb.Set(ctx, key, val, 0).Err()
		return err
	
	//Total number of movies for a given release year
	case 3:
		year := ScanStringWithPrompt("Write a release year to count the total number of movies: ", scanner)
		key := "total_movies_" + year
		releaseYear, err := strconv.Atoi(year)
		if err != nil {
		fmt.Println(year, "is not a valid release year")
		return nil
		}
		val, err := rdb.Get(ctx, key).Result()
		if err == nil {
			fmt.Println("Found in cache")
			fmt.Println(val)
			return nil
		}
		val = ""
		
		matchMovies := bson.D{primitive.E{Key: "$match", Value: bson.D{
			primitive.E{Key: "type", Value: "Movie"},
			primitive.E{Key: "release_year", Value: releaseYear},
		}}}
		groupMovies := bson.D{primitive.E{Key: "$group", Value: bson.D{
			primitive.E{Key: "_id", Value: "movies_year"},
			primitive.E{Key: "count", Value: bson.D{
				primitive.E{Key: "$sum", Value: 1}}},
		}}}
		
		cursor, err := GetTitlesColl(client).Aggregate(ctx, mongo.Pipeline{matchMovies, groupMovies})
		
		if err != nil {
			return err
		}
		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			return err
		}
		val += fmt.Sprintf("Movies in %d : ", releaseYear)
		for _, result := range results {
			val += fmt.Sprintf(" - %d\n", result["count"])
		}
		
		fmt.Println(val)
		err = rdb.Set(ctx, key, val, 0).Err()
		return err
	
	//Top 10 actors in more movies
	case 4:
		key := "top_10_actors"
		val, err := rdb.Get(ctx, key).Result()
		if err == nil {
			fmt.Println("Found in cache")
			fmt.Println(val)
			return nil
		}
		val = ""
		
		unwindCast := bson.D{primitive.E{Key: "$unwind", Value: "$cast"}}
		groupActors := bson.D{primitive.E{Key: "$group", Value: bson.D{
			primitive.E{Key: "_id", Value: "$cast"},
			primitive.E{Key: "count", Value: bson.D{
				primitive.E{Key: "$sum", Value: 1}}},
		}}}
		sort := bson.D{primitive.E{Key: "$sort", Value: bson.D{
		  primitive.E{Key: "count", Value: -1},
		}}}
		limit := bson.D{primitive.E{Key: "$limit", Value: 10}}
		
		cursor, err := GetTitlesColl(client).Aggregate(ctx, mongo.Pipeline{unwindCast, groupActors, sort, limit})
		
		if err != nil {
			return err
		}
		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			return err
		}
		val += fmt.Sprintln("Top 10 actors starring in most films")
		for _, result := range results {
			val += fmt.Sprintf(" - %s\n", result["_id"])
		}
		
		fmt.Println(val)
		err = rdb.Set(ctx, key, val, 0).Err()
		return err
		
	//Top 10 directors with most movies
	case 5:
		key := "top_10_directors"
		val, err := rdb.Get(ctx, key).Result()
		if err == nil {
			fmt.Println("Found in cache")
			fmt.Println(val)
			return nil
		}
		val = ""
		
		unwindDirector := bson.D{primitive.E{Key: "$unwind", Value: "$director"}}
		groupDirectors := bson.D{primitive.E{Key: "$group", Value: bson.D{
			primitive.E{Key: "_id", Value: "$director"},
			primitive.E{Key: "count", Value: bson.D{
				primitive.E{Key: "$sum", Value: 1}}},
		}}}
		sort := bson.D{primitive.E{Key: "$sort", Value: bson.D{
		  primitive.E{Key: "count", Value: -1},
		}}}
		limit := bson.D{primitive.E{Key: "$limit", Value: 10}}
		
		cursor, err := GetTitlesColl(client).Aggregate(ctx, mongo.Pipeline{unwindDirector, groupDirectors, sort, limit})
		
		if err != nil {
			return err
		}
		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			return err
		}
		val += fmt.Sprintln("Top 10 directors with most films")
		for _, result := range results {
			val += fmt.Sprintf(" - %s\n", result["_id"])
		}
		
		fmt.Println(val)
		err = rdb.Set(ctx, key, val, 0).Err()
		return err
		
	default:
		fmt.Println("Action couldn't be identified")
	}
	return nil
}

func PrintMenuInConsole () {
	fmt.Println("Statistics")
	fmt.Println()
	fmt.Println("	1) Total number of movies and TV shows")
	fmt.Println("	2) Total number of movies for a given country")
	fmt.Println("	3) Total number of TV shows for a given release year")
	fmt.Println("	4) Top 10 actors appearing in more movies")
	fmt.Println("	5) Top 10 directors with more movies")
}
