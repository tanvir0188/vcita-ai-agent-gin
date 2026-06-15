build:
	@go build -o bin/vcita ./cmd/server

run: build
	@./bin/vcita

dev:
	@go run ./cmd/server

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations -seq $(name)

migrate-up:
	@go run ./cmd/migrate up

migrate-down:
	@go run ./cmd/migrate down

clean:
	@rm -rf bin