package controllers

import (
	"net/http"

	"github.com/wangyaodream/gerty-goblog/pkg/view"
)

type CategoriesController struct {
	BaseController
}

func (*CategoriesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "categories.create")
}

func (*CategoriesController) Store(w http.ResponseWriter, r *http.Request) {

}
