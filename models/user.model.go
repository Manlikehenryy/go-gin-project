package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string               `json:"firstName" binding:"required"`
	LastName  string               `json:"lastName" binding:"required"`
	Email     string               `json:"email" binding:"required"`
	Password  []byte               `json:"-"`
	Phone     string               `json:"phone" binding:"required"`
}



func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
