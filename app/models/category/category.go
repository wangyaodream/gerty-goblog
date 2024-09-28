package category

import (
	"github.com/wangyaodream/gerty-goblog/app/models"
	"github.com/wangyaodream/gerty-goblog/pkg/route"
)

type Category struct {
	models.BaseModel

	Name string `gorm:"type:varchar(255);not null;" valid:"name"`
}

func (category Category) Link() string {
	return route.Name2URL("categories.show", "id", category.GetStringID())
}
