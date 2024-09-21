package middlewares

import (
	"net/http"

	"github.com/wangyaodream/gerty-goblog/pkg/auth"
	"github.com/wangyaodream/gerty-goblog/pkg/flash"
)

func Auth(next HTTPHandlerFunc) HTTPHandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        if !auth.Check() {
            flash.Warning("登录用户才能访问此页面")
            http.Redirect(w, r, "/", http.StatusFound)
            return
        }
        next(w, r)
    }
}
