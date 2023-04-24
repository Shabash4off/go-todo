package handlers

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"todo/internal/todo/models"
	"todo/internal/todo/repository"
)

type Todo struct {
	repository repository.TodoRepositoryInterface
}

func NewTodoHandler(repository repository.TodoRepositoryInterface) *Todo {
	return &Todo{repository: repository}
}

func (t *Todo) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "JSON decoding error", http.StatusBadRequest)
		return
	}

	result, err := t.repository.Create(todo)
	if err != nil {
		http.Error(w, "Todo creation error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result.InsertedID); err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func (t *Todo) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	todos, err := t.repository.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(w, "Todo reading error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func (t *Todo) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idParam := r.URL.Query().Get("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	todo, err := t.repository.GetByID(id)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func (t *Todo) UpdateByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idParam := r.URL.Query().Get("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	var todo models.Todo
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "JSON decoding error", http.StatusBadRequest)
		return
	}

	result, err := t.repository.UpdateByID(id, bson.M{
		"$set": bson.M{
			"title":   todo.Title,
			"content": todo.Content,
			"status":  todo.Status,
		},
	})
	if err != nil {
		http.Error(w, "Todo updating error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func (t *Todo) DeleteByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idParam := r.URL.Query().Get("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	result, err := t.repository.DeleteByID(id)
	if err != nil {
		http.Error(w, "Todo deleting error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}
