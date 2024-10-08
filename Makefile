include .env

DATABASE_URL = "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=require"
MIGRATIONS_PATH = "migrations"

# build the api binary
build:
	@go build -o bin/api cmd/main.go

# remove the api binary
clean:
	@rm -rf bin/api

# build and run the api binary
run: clean build
	@./bin/api

# install all dependencies
install:
	@go get -u ./...
	@go mod tidy

# run the tests
test:
	@go test -v ./...

# create a database migration file
migration-create:
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(filter-out $@,$(MAKECMDGOALS))

# run the database migrations
migration-up:
	@migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) up

# rollback the database migrations
migration-down:
	@migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) down

.PHONY: build clean run install test migration-create migration-up migration-down