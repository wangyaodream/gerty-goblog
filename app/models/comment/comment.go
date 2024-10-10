package comment

import "github.com/wangyaodream/gerty-goblog/app/models"

type Comment struct {
	models.BaseModel

	Content string `gorm:"type:longtext;not null;" valid:"content"`
	UserID  uint64 `gorm:"default:0;index"`
	PostID  uint64 `gorm:"not null;index"`
}
