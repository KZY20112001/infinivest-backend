package models

type User struct {
	ID           string `json:"id" bson:"_id"`
	Email        string `json:"email" bson:"email"`
	PasswordHash string `json:"password_hash" bson:"password_hash"`
}
