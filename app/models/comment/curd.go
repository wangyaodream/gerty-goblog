package comment

import (
	"github.com/wangyaodream/gerty-goblog/pkg/model"
	"github.com/wangyaodream/gerty-goblog/pkg/types"
)

func Get(idstr string) (Comment, error) {
	var comment Comment

	id := types.StringToUint64(idstr)
	return comment, model.DB.Preload("User").Preload("Article").First(&comment, id).Error
}
