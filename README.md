# URL Management 

A Gin framework based REST API server.

## Setup using Docker (Recommended)
### Prerequisites
* Docker daemon should be up and running

run command
```
docker-compose up -d
```
Server will be setup including db and runs on http://localhost:8080

## Manual setup 
### Prerequisites
* Setup postgres db, install and setup go

* create a .env file with below variables
    ```
    DB_HOST=localhost
    DB_USER=dd
    DB_PASSWORD=root
    DB_NAME=url-management
    DB_PORT=5432
    JWT_KEY='urlManagementSecretForAuthentication'
    ```

## Setup
Install Dependencies
```
go mod tidy
```

Run the server
```
go run main.go
```

Server will be up and running on http://localhost:8080
