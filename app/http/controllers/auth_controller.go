package controllers

import (
	"net/http"

	"github.com/wangyaodream/gerty-goblog/pkg/view"
)

type AuthController struct {
}

func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.Render(w, ArticlesFormData{}, "auth.register")
}

func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {

}
