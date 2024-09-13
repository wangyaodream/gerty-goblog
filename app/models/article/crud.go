package article

import (
	"github.com/wangyaodream/gerty-goblog/pkg/model"
	"github.com/wangyaodream/gerty-goblog/pkg/types"
)


// Get post from id
func Get(idstr string) (Article, error) {
    var article Article
    id := types.StringToUint64(idstr)
    if err := model.DB.First(&article, id).Error; err != nil {
        return article, err
    }

    return article, nil
}

// get all post
func GetAll() ([]Article, error) {
    var articles []Article
    if err := model.DB.Find(&articles).Error; err != nil {
        return articles, err
    }
    return articles, nil
}
