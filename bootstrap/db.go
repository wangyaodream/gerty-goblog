package bootstrap

import (
	"time"

	"github.com/wangyaodream/gerty-goblog/pkg/model"
)

func SetupDB() {

	// create database connection
	db := model.ConnectDB()

	// print the request of database with cli
	sqlDB, _ := db.DB()

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}
