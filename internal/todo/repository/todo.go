package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
	"todo/internal/todo/models"
)

type TodoRepositoryInterface interface {
	Create(todo models.Todo) (*mongo.InsertOneResult, error)
	GetAll() ([]*models.Todo, error)
	GetByID(id primitive.ObjectID) (*models.Todo, error)
	UpdateByID(id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error)
	DeleteByID(id primitive.ObjectID) (*mongo.DeleteResult, error)
}

const TodoCollection = "todos"

type TodoRepository struct {
	db *mongo.Database
}

func NewTodoRepository(db *mongo.Database) *TodoRepository {
	return &TodoRepository{db: db}
}

func (t *TodoRepository) Create(todo models.Todo) (*mongo.InsertOneResult, error) {
	collection := t.db.Collection(TodoCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return collection.InsertOne(ctx, todo)
}

func (t *TodoRepository) GetAll() ([]*models.Todo, error) {
	var todos []*models.Todo

	collection := t.db.Collection(TodoCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("Failed to close cursor: %v", err)
		}
	}()

	for cursor.Next(ctx) {
		var todo models.Todo
		err := cursor.Decode(&todo)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (t *TodoRepository) GetByID(id primitive.ObjectID) (*models.Todo, error) {
	var todo models.Todo

	collection := t.db.Collection(TodoCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&todo)
	return &todo, err
}

func (t *TodoRepository) UpdateByID(id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	collection := t.db.Collection(TodoCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return collection.UpdateByID(ctx, id, update)
}

func (t *TodoRepository) DeleteByID(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection := t.db.Collection(TodoCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return collection.DeleteOne(ctx, bson.M{"_id": id})
}
