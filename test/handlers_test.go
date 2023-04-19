package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"todo/internal/api/handlers"
	"todo/internal/todo/models"
)

type FakeTodoRepository struct {
	data map[primitive.ObjectID]*models.Todo
}

func NewFakeTodoRepository() *FakeTodoRepository {
	return &FakeTodoRepository{
		data: make(map[primitive.ObjectID]*models.Todo),
	}
}

func (f *FakeTodoRepository) Create(todo models.Todo) (*mongo.InsertOneResult, error) {
	todo.ID = primitive.NewObjectID()
	f.data[todo.ID] = &todo
	return &mongo.InsertOneResult{InsertedID: todo.ID}, nil
}

func (f *FakeTodoRepository) GetAll() ([]*models.Todo, error) {
	todos := make([]*models.Todo, 0, len(f.data))
	for _, todo := range f.data {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (f *FakeTodoRepository) GetByID(id primitive.ObjectID) (*models.Todo, error) {
	todo, ok := f.data[id]
	if !ok {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

func (f *FakeTodoRepository) UpdateByID(id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	todo, ok := f.data[id]
	if !ok {
		return nil, errors.New("todo not found")
	}

	if title, ok := update["title"].(string); ok {
		todo.Title = title
	}
	if content, ok := update["content"].(string); ok {
		todo.Content = content
	}
	if status, ok := update["status"].(models.Status); ok {
		todo.Status = status
	}

	return &mongo.UpdateResult{ModifiedCount: 1}, nil
}

func slicesContains(n models.Todo, ns []models.Todo) bool {
	for _, todo := range ns {
		if n.ID.Hex() == todo.ID.Hex() {
			return true
		}
	}
	return false
}

func slicesEqual(a, b []models.Todo) bool {
	if len(a) != len(b) {
		return false
	}

	for _, todo := range a {
		if !slicesContains(todo, b) {
			return false
		}
	}
	return true
}

func (f *FakeTodoRepository) DeleteByID(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	if _, ok := f.data[id]; !ok {
		return nil, errors.New("todo not found")
	}
	delete(f.data, id)
	return &mongo.DeleteResult{DeletedCount: 1}, nil
}

func TestTodoCreate(t *testing.T) {
	fakeRepository := NewFakeTodoRepository()

	handler := handlers.NewTodoHandler(fakeRepository)

	todo := models.Todo{
		Title:   "Test",
		Content: "Test content",
		Status:  models.TodoStatusPending,
	}

	todoJSON, err := json.Marshal(todo)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/todo/create", bytes.NewReader(todoJSON))
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	w := httptest.NewRecorder()

	handler.Create(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestTodoGetAll(t *testing.T) {
	fakeRepository := NewFakeTodoRepository()

	todos := make([]models.Todo, 0, 3)
	for i := 1; i <= 3; i++ {
		todo := models.Todo{
			ID:      primitive.NewObjectID(),
			Title:   fmt.Sprintf("Test %v", i),
			Content: fmt.Sprintf("Test %v", i),
		}

		fakeRepository.data[todo.ID] = &todo
		todos = append(todos, todo)
	}

	handler := handlers.NewTodoHandler(fakeRepository)

	req, err := http.NewRequest(http.MethodGet, "/todos", nil)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var result []models.Todo
	err = json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if !slicesEqual(result, todos) {
		t.Errorf("Handler returned wrong body: got %v want %v", result, todos)
	}
}

func TestTodoGetByID(t *testing.T) {
	fakeRepository := NewFakeTodoRepository()

	todo := models.Todo{
		ID:      primitive.NewObjectID(),
		Title:   "Test",
		Content: "Testing",
	}

	fakeRepository.data[todo.ID] = &todo

	handler := handlers.NewTodoHandler(fakeRepository)

	req, err := http.NewRequest(http.MethodGet, "/todo/id", nil)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	q := req.URL.Query()
	q.Add("id", todo.ID.Hex())
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var result models.Todo
	err = json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if !reflect.DeepEqual(result, todo) {
		t.Errorf("Handler returned wrong body: got %v want %v", result, todo)
	}
}

func TestTodoUpdateByID(t *testing.T) {
	fakeRepository := NewFakeTodoRepository()

	todo := models.Todo{
		ID:      primitive.NewObjectID(),
		Title:   "Test",
		Content: "Testing",
	}

	fakeRepository.data[todo.ID] = &todo

	update := bson.M{
		"title":   "Updated test",
		"content": "Updated testing",
		"status":  models.TodoStatusComplete,
	}

	updateJSON, err := json.Marshal(update)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	handler := handlers.NewTodoHandler(fakeRepository)

	req, err := http.NewRequest(http.MethodPut, "/todo/update", bytes.NewReader(updateJSON))
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	q := req.URL.Query()
	q.Add("id", todo.ID.Hex())
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	handler.UpdateByID(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var result mongo.UpdateResult
	err = json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if result.ModifiedCount != 1 {
		t.Errorf("Expected modified count to be 1, got %d", result.ModifiedCount)
	}
}

func TestTodoDeleteByID(t *testing.T) {
	fakeRepository := NewFakeTodoRepository()

	todo := models.Todo{
		ID:      primitive.NewObjectID(),
		Title:   "Test",
		Content: "Testing",
	}

	fakeRepository.data[todo.ID] = &todo

	handler := handlers.NewTodoHandler(fakeRepository)

	req, err := http.NewRequest(http.MethodDelete, "/todo/update", nil)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	q := req.URL.Query()
	q.Add("id", todo.ID.Hex())
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	handler.DeleteByID(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var result mongo.DeleteResult
	err = json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if result.DeletedCount != 1 {
		t.Errorf("Expected deleted count to be 1, got %d", result.DeletedCount)
	}
}
