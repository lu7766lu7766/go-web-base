package database

import (
	"fmt"
	"p4_web/model"
	"p4_web/tools/env"
	"p4_web/tools/exception"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dsn string

func init() {

	dbUser := env.Get("DB_USER")
	dbPwd := env.Get("DB_PASSWORD")
	dbHost := env.Get("DB_HOST")
	dbPort := env.Get("DB_PORT")
	dbName := env.Get("DB_NAME")
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPwd, dbHost, dbPort, dbName)

	// log.Fatal(con)
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		QueryFields: true,
	})

	if err != nil {
		// log.Print("failed to connect database")
		panic(exception.ApiException{
			Code:    []int{-3},
			Message: fmt.Sprintf("db connect error %v", err.Error()),
		})
	}

	if DB.Error != nil {
		// log.Printf("database error %v", DB.Error)
		panic(exception.ApiException{
			Code:    []int{-3},
			Message: fmt.Sprintf("database error %v", DB.Error),
		})
	}

	DB.AutoMigrate(&model.User{})
}
