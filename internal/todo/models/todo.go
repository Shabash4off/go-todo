package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Status uint8

const (
	TodoStatusPending Status = iota
	TodoStatusComplete
)

type Todo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	Status    Status             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
