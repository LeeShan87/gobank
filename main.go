package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	listenAddr := os.Getenv("LISTEN_ADDR")
	dbConfig := &PostgresStoreConfig{
		user:     os.Getenv("POSTGRES_USER"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		port:     os.Getenv("POSTGRES_PORT"),
		dbName:   os.Getenv("POSTGRES_DB"),
	}
	storage, err := NewPostgressStore(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	if err := storage.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", storage)
	server := NewApiServer(fmt.Sprintf(":%s", listenAddr), storage)
	server.Run()
}
