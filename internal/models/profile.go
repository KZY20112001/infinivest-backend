package models

type Profile struct {
	UserID    string `json:"user_id" bson:"user_id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
}
