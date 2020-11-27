package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx, client, cancel, err := ConnectMongoDB()
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       5,
	})
	rdb.ConfigSet(ctx, "maxmemory", "2mb")
	rdb.ConfigSet(ctx, "maxmemory-policy", "allkeys-lru")

	if err = InitMongoDB(ctx, client); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for action := 0; action != 5; {
		fmt.Println()
		fmt.Println("	1) Search for a movie/actor/TV show")
		fmt.Println("	2) Get statistics for movies/TV shows")
		fmt.Println("	3) Add a new movie/TV show")
		fmt.Println("	4) Update a movie/TV show")
		fmt.Println("	5) Exit")
		actionStr := ScanStringWithPrompt("Select an action: ", scanner)
		action, err = strconv.Atoi(actionStr)
		if err != nil {
			action = 0
		}
		switch action {
		case 1:
			if err = Query(ctx, client, scanner, rdb); err != nil {
				log.Fatal(err)
			}
		case 2:
			if err = Statistics(ctx, client, scanner, rdb); err != nil {
				log.Fatal(err)
			}
		case 3:
			if err = Add(ctx, client, scanner, rdb); err != nil {
				log.Fatal(err)
			}
		case 4:
			if err = Update(ctx, client, scanner, rdb); err != nil {
				log.Fatal(err)
			}
		case 5:
			fmt.Println("Bye!")
		default:
			fmt.Println("Action couldn't be parsed")
		}
	}
}
