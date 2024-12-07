package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	UserID    uint
	User      User
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
}
