package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title  string             `json:"title" bson:"title" binding:"required"`
	Desc   string             `json:"desc" bson:"desc" binding:"required"`
	Status bool               `json:"status" bson:"status"`
	UserId primitive.ObjectID `json:"userId" bson:"userId"`
	// User   User               `json:"user" bson:"user,omitempty"`
}
