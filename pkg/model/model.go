package model

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/wangyaodream/gerty-goblog/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	godotenv.Load(".env")
	var err error

	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/goblog?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)
	postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=goblog port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
	)

	switch os.Getenv("DB_TYPE") {
	case "mysql":
		config := mysql.New(mysql.Config{
			DSN: mysqlDSN,
		})
		DB, err = gorm.Open(config, &gorm.Config{})
	case "postgres":
		config := postgres.New(postgres.Config{
			DSN: postgresDSN,
		})
		DB, err = gorm.Open(config, &gorm.Config{})
	}

	logger.LogError(err)

	return DB
}
