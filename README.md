# Blog API

A simple RESTful API for a blog platform implemented in Go.

## Features

- CRUD operations for blog posts (Create, Read, Update, Delete)
- In-memory storage with JSON seed data
- Pagination support on GET /posts endpoint
- Request logging middleware using `zerolog`
- Unit tests with mocks for handlers and services

## Project Structure
```
├── cmd/server # Main server entrypoint
├── internal
│ ├── handler # HTTP handlers, middleware, routes, mocks
│ ├── model # Data models (Post, NewPost)
│ ├── service # Business logic
│ └── storage # In-memory data storage and seed data
```

## Getting Started

### Prerequisites

- Go 1.24 

### Run the server

```bash
go run ./cmd/server
```

The server will start on port 8080.

### API Endpoints
| Method | Endpoint          | Description                                           |
|--------|-------------------|-----------------------------------------------------|
| GET    | `/posts`          | List blog posts (supports `limit` and `offset` query params) |
| GET    | `/posts/{id}`     | Get a single post by ID                              |
| POST   | `/posts`          | Create a new post                                   |
| PUT    | `/posts/{id}`     | Update a post by ID                                 |
| DELETE | `/posts/{id}`     | Delete a post by ID                                 |

### Run tests

```bash
go test ./internal/...
```

## Notes
- Uses in-memory storage, so all data is lost on restart.
- Seed data is loaded from `internal/storage/seed-posts/posts.json` by default.
- Mocks for unit tests are generated with `testify/mock` and located in `internal/handler/mocks` and `internal/service/mocks`.
