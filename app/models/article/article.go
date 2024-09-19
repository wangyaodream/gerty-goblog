package article

import (
	"strconv"

	"github.com/wangyaodream/gerty-goblog/app/models"
	"github.com/wangyaodream/gerty-goblog/pkg/route"
)

type Article struct {
	// ID    uint64
	models.BaseModel
	Title string `gorm:"type:varchar(255);not null;" valid:"title"`
	Body  string `gorm:"type:longtext;not null;" valid:"body"`
}

func (a Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatUint(a.ID, 10))
}
