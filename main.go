package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func seedAccount(db Storage, fName, lName, passWD string) {
	acc, err := NewAccount(fName, lName, passWD)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}
}

func seedDB(db Storage) {
	fmt.Println("Populating database with seed data")
	seedAccount(db, "Anthony", "GG", "hunter88888")
	seedAccount(db, "Bob", "Russ", "rusty_999")
	seedAccount(db, "Zoltan", "Toma", "it_is_secret")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	seed := flag.Bool("seed", false, "Seed the database")
	flag.Parse()
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
	if *seed {
		fmt.Println("Seeding the database")
		if err := storage.dropTableForSeed(); err != nil {
			log.Fatal(err)
		}
		if err := storage.Init(); err != nil {
			log.Fatal(err)
		}
		seedDB(storage)
		fmt.Println("Seeding finished")
		os.Exit(0)
	}
	if err := storage.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", storage)
	server := NewApiServer(fmt.Sprintf(":%s", listenAddr), storage)
	server.Run()
}
