package article

import (
	"net/http"

	"github.com/wangyaodream/gerty-goblog/pkg/logger"
	"github.com/wangyaodream/gerty-goblog/pkg/model"
	"github.com/wangyaodream/gerty-goblog/pkg/pagination"
	"github.com/wangyaodream/gerty-goblog/pkg/route"
	"github.com/wangyaodream/gerty-goblog/pkg/types"
)

// Get post from id
func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToUint64(idstr)
	if err := model.DB.Preload("User").First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

// get all post
func GetAll(r *http.Request, perPage int) ([]Article, pagination.ViewData, error) {
	db := model.DB.Model(Article{}).Order("created_at desc")
	_pager := pagination.New(r, db, route.Name2URL("home"), perPage)

	viewData := _pager.Paging()

	var articles []Article
	_pager.Results(&articles)

	return articles, viewData, nil
}

// 获取全部文章
func GetByUserID(uid string) ([]Article, error) {
	var articles []Article
	if err := model.DB.Where("user_id = ?", uid).Preload("User").Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

// create post by aritcle.ID
func (article *Article) Create() (err error) {
	result := model.DB.Create(&article)
	if err = result.Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

func (article *Article) Update() (rowsAffected int64, err error) {
	result := model.DB.Save(&article)
	if err = result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}

	return result.RowsAffected, nil
}

func (article *Article) Delete() (rowsAffected int64, err error) {
	result := model.DB.Delete(&article)
	if err = result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}

	return result.RowsAffected, nil
}

func GetByCategoryID(cid string, r *http.Request, perPage int) ([]Article, pagination.ViewData, error) {
	// 初始化分页
	db := model.DB.Model(Article{}).Where("category_id = ?", cid).Order("created_at desc")
	_pager := pagination.New(r, db, route.Name2URL("categories.show", "id", cid), perPage)

	// 获取视图数据
	viewData := _pager.Paging()

	// 获取文章数据
	var articles []Article
	_pager.Results(&articles)

	return articles, viewData, nil
}
