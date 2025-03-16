package models

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	UserID     uint
	User       User
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Address    string `json:"address"`
	ProfileUrl string `json:"profileUrl"`
	ProfileID  string `json:"profileID"`
}
