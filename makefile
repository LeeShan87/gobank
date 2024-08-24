build:
	@go build -o bin/gobank

run: build
	@./bin/gobank

test:
	@go test -v ./...

clean:
	rm -rf ./bin