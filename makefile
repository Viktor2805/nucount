build:
	go build cmd/timesaver/main.go && ./main 
run: 
	go run cmd/timesaver/main.go
migration_up:
	migrate -path pkg/db/migrations/ -database "postgresql://postgres:1@localhost:5435/time-saver?sslmode=disable" -verbose up

migration_down:
	migrate -path pkg/db/migrations/ -database "postgresql://postgres:1@localhost:5435/time-saver?sslmode=disable" -verbose down

migration_fix:
	migrate -path pkg/db/migrations/ -database "postgresql://postgres:1@localhost:5435/time-saver?sslmode=disable" force 2

migration_create:
	migrate create -ext sql -dir pkg/db/migrations/ -seq create_time_entry_table
	
