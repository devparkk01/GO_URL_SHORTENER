
# URL Shortener API

A simple URL shortening API built using Go, Gorilla Mux, and SQLite. This project provides an API for creating short URLs and retrieving the original URLs.

## Features

- Create short URLs
- Retrieve original URLs
- Delete short URLs
- Update short URLs
- Support for concurrent requests

## Technologies Used

- Go (Golang)
- Gorilla Mux for routing
- SQLite for the database

## Getting Started

### Prerequisites

- Go 1.16 or later
- SQLite

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/devparkk01/GO_URL_SHORTENER.git
   ```

2. Change directory to the project folder:
   ```bash
   cd GO_URL_SHORTENER
   ```

3. Install Dependencies
   ```bash
   go mod tidy

   ```

4. Run the application:
   ```bash
   go run main.go
   ```
   The server will start on  http://localhost:8080.

### API Endpoints

- **POST /api/short**: Create a new short URL
  - POST http://localhost:8080/api/short 
  - Request Body 
  ```
    "original_url": "https://youtube.com/llkl79/abc"
  ```
  - Sample Response 
  ```
    "original_url": "https://youtube.com/llkl79/abc",
    "short_url": "28b6NWjU",
    "created_at": "2024-10-16 23:05:18"
  ```
- **GET /api/short/{shortUrl}**: Retrieve the original URL
    - GET http://localhost:8080/api/short/28b6NWjU
    - Sample Response
    ```
      "original_url": "https://youtube.com/llkl79/abc",
      "short_url": "28b6NWjU",
      "created_at": "2024-10-16 23:05:18"
    ```
- **PUT /api/short/{shortUrl}**: Update a short URL
  - PUT http://localhost:8080/api/short/28b6NWjU
  - Sample Response
  ```
    "updated_short_url": "i5oBH2ft"
  ```
- **DELETE /api/short/{shortUrl}**: Delete a short URL
    - DELETE http://localhost:8080/api/short/i5oBH2ft


### Running Tests

To run the tests, execute:
```bash
go test ./...
```

