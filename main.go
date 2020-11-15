package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NetflixTitle : Struct representing either a Movie or a TV Show
type NetflixTitle struct {
	Cast        []string `json:"cast"`
	Country     []string `json:"country"`
	DateAdded   string   `json:"date_added"`
	Description string   `json:"description"`
	Director    []string `json:"director"`
	Duration    string   `json:"duration"`
	ListedIn    []string `json:"listed_in"`
	Rating      string   `json:"rating"`
	ReleaseYear int      `json:"release_year"`
	ShowID      int      `json:"show_id"`
	Title       string   `json:"title"`
	TitleType   string   `json:"type"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	err = initDB(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
}
