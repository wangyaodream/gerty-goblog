package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/wangyaodream/gerty-goblog/pkg/logger"
)

var DB *sql.DB

func Initialize() {
	initDB()
	createTables()
}

func initDB() {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	dbType := os.Getenv("DB_TYPE")
	mysqlConnStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/goblog",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"))

	// PostgreSQL连接字符串
	postgresConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	switch dbType {
	case "mysql":
		DB, err = connectDB("mysql", mysqlConnStr)
		logger.LogError(err)

		// 配置连接属性
		DB.SetMaxOpenConns(25)                 // 最大连接数
		DB.SetMaxIdleConns(25)                 // 最大空闲数
		DB.SetConnMaxLifetime(5 * time.Minute) // 每个链接的过期时间
	case "postgres":
		DB, err = connectDB("postgres", postgresConnStr)
		logger.LogError(err)
	}
}

func connectDB(dbType, connStr string) (*sql.DB, error) {
	db, err := sql.Open(dbType, connStr)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

func createTables() {
	dbType := os.Getenv("DB_TYPE")
	var createArticlesSQL string
	switch dbType {
	case "mysql":
		createArticlesSQL = `CREATE TABLE IF NOT EXISTS articles(
        id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
        title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
        body longtext COLLATE utf8mb4_unicode_ci
    ); `
	case "postgres":
		createArticlesSQL = `CREATE TABLE IF NOT EXISTS articles (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        body TEXT
    );`

	}

	_, err := DB.Exec(createArticlesSQL)
	logger.LogError(err)
}
