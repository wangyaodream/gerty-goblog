package comment

import (
	"github.com/wangyaodream/gerty-goblog/app/models"
	"github.com/wangyaodream/gerty-goblog/app/models/article"
	"github.com/wangyaodream/gerty-goblog/app/models/user"
)

type Comment struct {
	models.BaseModel

	Content   string `gorm:"type:longtext;not null;" valid:"content"`
	UserID    uint64 `gorm:"not null;index"`
	User      user.User
	ArticleID uint64 `gorm:"not null;index"`
	Article   article.Article
}
