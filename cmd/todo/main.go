package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"todo/internal/api/handlers"
	"todo/internal/db"
	"todo/internal/todo/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	mongoUri, ok := os.LookupEnv("MONGO_URI")
	if !ok {
		log.Fatal("MONGO_URI not found")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatal("PORT not found")
	}

	host, _ := os.LookupEnv("HOST")
	address := host + ":" + port

	mongoClient := db.NewMongoClient(mongoUri)
	todoDb := mongoClient.Database("todo")

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Failed to disconect from mongodb: %v", err)
		}
	}()

	todoRepository := repository.NewTodoRepository(todoDb)
	todoHandler := handlers.NewTodoHandler(todoRepository)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/todo/create/", todoHandler.Create)
	mux.HandleFunc("/api/todos/", todoHandler.GetAll)
	mux.HandleFunc("/api/todo/id/", todoHandler.GetByID)
	mux.HandleFunc("/api/todo/update/", todoHandler.UpdateByID)
	mux.HandleFunc("/api/todo/delete/", todoHandler.DeleteByID)

	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(address, mux))
}
