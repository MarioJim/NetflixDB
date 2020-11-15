package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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
	ctx, client, cancel, err := ConnectMongoDB()
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	err = InitMongoDB(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for action := 0; action != 5; {
		fmt.Println("	1) Search for a movie/actor/TV show")
		fmt.Println("	2) Get statistics for movies/TV shows")
		fmt.Println("	3) Add a new movie/TV show")
		fmt.Println("	4) Update a movie/TV show")
		fmt.Println("	5) Exit")
		fmt.Print("Select an action: ")
		scanner.Scan()
		action, err = strconv.Atoi(scanner.Text())
		if err != nil {
			action = 0
		}
		switch action {
		case 1:
			Query()
			break
		case 2:
			Statistics()
			break
		case 3:
			Add()
			break
		case 4:
			Update()
			break
		case 5:
			fmt.Println("Bye!")
			break
		default:
			fmt.Println("Action couldn't be parsed")
		}
	}
}
