package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	First_name string
	Last_name  string
	Email      string
	Username   string
	Password   string
	IsAdmin    string `gorm:"default:no;type:varchar(10)"`
}

type Errors struct {
	Errors string
}
