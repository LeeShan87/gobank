package main

func main() {
	server := NewApiServer(":4000")
	server.Run()
}
