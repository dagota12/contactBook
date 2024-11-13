package domain



type User struct {
	Phone string `json:"phone" bson:"phone"`
	Username string `json:"username" bson:"username"`
}