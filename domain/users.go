package domain

import "go.mongodb.org/mongo-driver/bson/primitive"



type User struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Phone string `json:"phone" bson:"phone"`
	Username string `json:"username" bson:"username"`
}