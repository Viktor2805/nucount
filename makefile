# Makefile for a Dna Analyzer project

#Building the application
build:
	@echo "Building..."

	go build cmd/main.go  

#Running the linter for application
lint:
	golangci-lint run ./...

# Test only internal packages
test:
	go test ./internal/... -v

# Test the entire module
test-all:
	go test ./... -v

# Run benchmarks
bench:
	go test -bench . -benchmem ./...

.PHONY: swagger
swagger:
	@echo "Generating Swagger..."
	swag init -g ./cmd/main.go -o ./docs

#Running the application
run: 
	go run cmd/main.go
	
# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

docker-build:
	docker build -t viktor2805/dna-analyzer:1.1 . && docker push viktor2805/dna-analyzer:1.1
	
# migration_up:
# 	migrate -path pkg/db/migrations/ -database "postgresql://postgres:1@localhost:5435/time-saver?sslmode=disable" -verbose up

# migration_down:
# 	migrate -path pkg/db/migrations/ -database "postgresql://postgres:1@localhost:5435/time-saver?sslmode=disable" -verbose down

# migration_fix:
# 	migrate -path pkg/db/migrations/ -database "postgresql://postgres:1@localhost:5435/time-saver?sslmode=disable" force 2

# migration_create:
# 	migrate create -ext sql -dir pkg/db/migrations/ -seq create_time_entry_table
	