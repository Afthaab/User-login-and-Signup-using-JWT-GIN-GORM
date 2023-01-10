package config

import (
	"fmt"

	"github.com/project_login/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "projectone"
)

func DBConn() (DB *gorm.DB) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	DB, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		fmt.Println("the error is here")
		panic(err)

	}

	DB.AutoMigrate(&models.User{})
	
	return DB
}
