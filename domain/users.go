package domain

import "go.mongodb.org/mongo-driver/bson/primitive"



type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string `json:"username" bson:"username,omitempty"`
	Phone     string `json:"phone" bson:"phone,omitempty"`
}

type  UpdateReq struct {
	Phone       string `json:"phone"`
	Username    string `json:"username"`
	NewUsername string `json:"newUsername,omitempty"`
	NewPhone    string `json:"newPhone,omitempty"`
}