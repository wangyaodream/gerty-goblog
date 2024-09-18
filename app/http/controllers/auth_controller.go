package controllers

import (
	"fmt"
	"net/http"

	"github.com/wangyaodream/gerty-goblog/app/models/user"
	"github.com/wangyaodream/gerty-goblog/app/requests"
	"github.com/wangyaodream/gerty-goblog/pkg/view"
)

type AuthController struct {
}

func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {

	// 获取表单数据
	_user := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}

	// 获取表单规则
	errs := requests.ValidateRegistrationForm(_user)

	if len(errs) > 0 {
		// 发生错误输出错误信息
		view.RenderSimple(w, view.D{
			"Errors": errs,
			"User":   _user,
		}, "auth.register")
	} else {
		_user.Create()

		if _user.ID > 0 {
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "注册失败!")
		}
	}

}

// Login
func (*AuthController) Login(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.login")
}

func (*AuthController) DoLogin(w http.ResponseWriter, r *http.Request) {
	// TODO 登录验证
}
