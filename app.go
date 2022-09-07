package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

var (
	db *redis.Client
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Server started.")
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"), nil))
}
