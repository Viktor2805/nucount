# DNA nucleotides count Service

This project contains a Go-based transactions service which uses PostgreSQL as its database. The service is containerized using Docker and Docker Compose.

## Prerequisites

- Docker
- Docker Compose

## Getting Started

### Building and Running the Application

1. Clone the repository to your local machine.
2. Navigate to the project directory.
3. Run the following command to build and start the services:

    ```sh
    docker-compose up --build
    ```

### Accessing the Application

Once the services are up and running, you can access the application and its API documentation as follows:

- **Transactions Service**: `http://localhost:3000`
- **Swagger API Documentation**: `http://localhost:3000/swagger/index.html`

### Running Tests

To run the tests for the entire project, use the following command:

```sh
docker-compose run --rm tests
 ```
 
### Environment Variables

The `docker-compose.yml` file contains the following environment variables for the `transactions` service:

- `POSTGRES_USER`: The username for the PostgreSQL database (default: `postgres`).
- `POSTGRES_PASSWORD`: The password for the PostgreSQL database (default: `1`).
- `POSTGRES_DB`: The name of the PostgreSQL database (default: `transactions`).
- `POSTGRES_HOST`: The hostname for the PostgreSQL service (default: `postgres`).
- `POSTGRES_PORT`: The port number for the PostgreSQL service (default: `5432`).
- `MIGRATION_URL`: The URL for database migrations (default: `file://pkg/migrations`).

### Accessing the PostgreSQL Database

The PostgreSQL database can be accessed at:

- Host: `localhost`
- Port: `5435`
- Username: `postgres`
- Password: `1`
- Database: `transactions`
## License

This project is licensed under the MIT License.
