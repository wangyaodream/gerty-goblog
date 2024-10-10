package bootstrap

import (
	"time"

	"github.com/wangyaodream/gerty-goblog/app/models/article"
	"github.com/wangyaodream/gerty-goblog/app/models/category"
	"github.com/wangyaodream/gerty-goblog/app/models/comment"
	"github.com/wangyaodream/gerty-goblog/app/models/user"
	"github.com/wangyaodream/gerty-goblog/pkg/model"
	"gorm.io/gorm"
)

func SetupDB() {

	// create database connection
	db := model.ConnectDB()

	// print the request of database with cli
	sqlDB, _ := db.DB()

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	migration(db)
}

func migration(db *gorm.DB) {
	// 注册自动迁移
	db.AutoMigrate(
		&user.User{},
		&article.Article{},
		&category.Category{},
		&comment.Comment{},
	)
}
