# URL Management 

A Gin framework based backend server.

## Prerequisites
* Setup postgres db

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