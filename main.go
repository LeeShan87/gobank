package main

import (
	"fmt"
	"log"
)

func main() {
	storage, err := NewPostgressStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := storage.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", storage)
	server := NewApiServer(":4000", storage)
	server.Run()
}
