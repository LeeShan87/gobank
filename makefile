build:
	@go build -o bin/gobank

run: build
	@./bin/gobank

test:
	@go test -v ./...

clean:
	rm -rf ./bin

up:
	docker compose up -d

down:
	docker compose down

seed: build
	@./bin/gobank --seed=true