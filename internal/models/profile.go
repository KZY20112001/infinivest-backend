package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	UserID     uint
	User       User
	FirstName  string
	LastName   string
	Address    string
	ProfileUrl string
	ProfileID  string
}
