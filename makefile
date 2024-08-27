build:
	@go build -o bin/gobank

run: build
	@./bin/gobank

test:
	@go test -v ./... --cover -count=1

test-html:
	@go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

clean:
	rm -rf ./bin

up:
	docker compose up -d

down:
	docker compose down

seed: build
	@./bin/gobank --seed=true