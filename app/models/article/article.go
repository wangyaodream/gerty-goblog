package article

import (
	"strconv"

	"github.com/wangyaodream/gerty-goblog/app/models"
	"github.com/wangyaodream/gerty-goblog/app/models/category"
	"github.com/wangyaodream/gerty-goblog/app/models/user"
	"github.com/wangyaodream/gerty-goblog/pkg/route"
	"github.com/wangyaodream/gerty-goblog/pkg/types"
)

type Article struct {
	// ID    uint64
	models.BaseModel
	Title      string `gorm:"type:varchar(255);not null;" valid:"title"`
	Body       string `gorm:"type:longtext;not null;" valid:"body"`
	UserID     uint64 `gorm:"not null;index"`
	User       user.User
	CategoryID uint64 `gorm:"not null;default:4;index"`
}

func (a Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatUint(a.ID, 10))
}

func (a Article) CreatedAtDate() string {
	return a.CreatedAt.Format("2006-01-02")
}

func (a Article) GetCategory() category.Category {
	c, _ := category.Get(types.Uint64ToString(a.CategoryID))
	return c
}
