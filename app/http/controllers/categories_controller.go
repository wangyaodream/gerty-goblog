package controllers

import (
	"fmt"
	"net/http"

	"github.com/wangyaodream/gerty-goblog/app/models/category"
	"github.com/wangyaodream/gerty-goblog/app/requests"
	"github.com/wangyaodream/gerty-goblog/pkg/view"
)

type CategoriesController struct {
	BaseController
}

func (*CategoriesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "categories.create")
}

func (*CategoriesController) Store(w http.ResponseWriter, r *http.Request) {
	_category := category.Category{
		Name: r.PostFormValue("name"),
	}

	errors := requests.ValidateCategoryForm(_category)

	if len(errors) == 0 {
		_category.Create()
		if _category.ID > 0 {
			fmt.Fprint(w, "创建成功！")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章分类失败")
		}
	} else {
		view.Render(w, view.D{
			"Category": _category,
			"Errors":   errors,
		}, "categories.create")
	}

}
