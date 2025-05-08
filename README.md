# ASAPP Chat Backend Challenge v1
### Overview
This is a Go based boilerplate which runs an HTTP Server configured to answer the endpoints defined in 
[the challenge you received](https://backend-challenge.asapp.engineering/).
All endpoints are configured in cmd/server.go and if you go deeper to the handlers
for each route, you will find a *TODO* comments where you are free to implement your solution.

### Instructions

They are located in the *docs/index.html* file

### Prerequisites

Installed Go version >= 1.12 since it uses Go Modules.

### How to run it

In project root:

`
go run cmd/server.go
`

### Run tests

In project root:

`
go test ./...
`

### Run it with docker locally

In project root: 

```bash
docker build -t messaging-app .
docker run -p 8080:8080 \
  -e JWT_SECRET_KEY=your-secret \
  messaging-app
```

## Endpoints

### Health Check

- **POST** `/check`  
  Returns a basic health status of the service and its database connection.

### Users

#### Create User

- **POST** `/users`
- **Request Body**:
  ```json
  {
    "username": "username",
    "password": "password"
  }
  ```
- **Response**:
  ```json
  {
    "id": 1
  }
  ```

### Authentication

#### Login

- **POST** `/login`
- **Request Body**:
  ```json
  {
    "username": "username",
    "password": "password"
  }
  ```
- **Response**:
  ```json
  {
    "id": 1,
    "token": "JWT_TOKEN_HERE"
  }
  ```

---

### Messages (Protected)

All `/messages` endpoints require a valid JWT token in the `Authorization` header:

```
Authorization: Bearer <token>
```

#### Get Messages

- **GET** `/messages`
- **Query Parameters**:
    - `start`: message ID to start from
    - `limit`: max number of messages to return
- **Response**:
  ```json
  [
    {
      "id": 1,
      "sender": 1,
      "recipient": 2,
      "timestamp": "2025-01-01T12:00:00Z",
      "content": {
        "type": "text",
        "text": "Hello"
      }
    }
  ]
  ```

#### Send Message

- **POST** `/messages`
- **Request Body**:
  ```json
  {
    "recipient": 2,
    "timestamp": "2025-01-01T12:00:00Z",
    "content": {
      "type": "text",
      "text": "Helo"
    }
  }
  ```
- **Response**:
  ```json
  {
    "id": 10
  }
  ```

## Environment Variables

| Variable         | Description                                                                |
|------------------|----------------------------------------------------------------------------|
| `JWT_SECRET_KEY` | Secret used to sign JWT tokens (required)                                  |
| `SQLITE_DSN`     | Path/DSN to the SQLite database file (Defaults to DB in memory if not set) |
