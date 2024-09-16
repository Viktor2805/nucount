# Makefile for a Dna Analyzer project

#Building the application
build:
	@echo "Building..."

	go build cmd/main.go  

#Running the linter for application
lint:
	revive ./...

# Test the application
test:
	go test ./... -v

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
	