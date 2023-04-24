# Todo Application  
This is a simple Todo application built using Go. The application allows users to create, read, update, and delete todo items. It uses MongoDB as the storage backend, and is structured with separate packages for handling HTTP requests, processing business logic, and interacting with the database.
## Table of Contents  
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Running the application](#running-the-application)
- [Using docker](#using-docker)
- [API Endpoints](#api-endpoints)
## Features
- Create a new todo item
- Retrieve all todo items
- Retrieve a specific todo item by ID
- Update a todo item by ID
- Delete a todo item by ID
## Requirements
- [Go](https://golang.org/doc/install) version 1.18 or higher
- [MongoDB](https://www.mongodb.com/try/download/community) version 4.4 or higher
## Installation
1. Clone the repository:
    ```bash
    git clone https://github.com/Shabash4off/go-todo.git
    ```
2. Change to the project directory:
    ```bash
    cd go-todo
    ```
3. Download the required Go modules:
    ```bash
    go mod download
    ```
## Running the application
1. Ensure that MongoDB is running on your system.
2. Create .env and fill it
   ```bash
   cp .env.sample .env
   ```
3. Build and run the application:
    ```bash
    go build -o todo todo/cmd/todo
    ./todo 
    ```
## Using Docker
You can also run the application using Docker and Docker Compose.
### Prerequisites
- Docker
- Docker Compose
### Build and Run the Application
1. Build the Docker image:
   ```bash
   docker build -t todo .
   ```
2. Run the Docker container:
   ```bash
   docker run -p 8080:8080 todo
   ```
Now the application is running at `http://localhost:8080`.
### Using Docker Compose
1. Start the application and MongoDB using Docker Compose:
   ```bash
   docker-compose up
   ```
2. Stop the application and MongoDB:
   ```bash
   docker-compose down
   ```

## API Endpoints
| Method   | Endpoint                  | Description                |
|----------|---------------------------|----------------------------|
| POST     | /api/todo/create          | Create a new todo item     |
| GET	     | /api/todos                | Retrieve all todo items    |
| GET	     | /api/todo/id?id={id}      | Retrieve a todo item by ID |
| PUT      | /api/todo/update?id={id}  | Update a todo item by ID   |
| DELETE   | /api/todo/delete?id={id}	 | Delete a todo item by ID   |
Note: Replace {id} with the actual ID of the todo item in the URL when using the specific item endpoints.