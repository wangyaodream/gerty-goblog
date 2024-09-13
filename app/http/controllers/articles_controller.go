package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/wangyaodream/gerty-goblog/app/models/article"
	"github.com/wangyaodream/gerty-goblog/pkg/logger"
	"github.com/wangyaodream/gerty-goblog/pkg/route"
	"github.com/wangyaodream/gerty-goblog/pkg/types"
	"gorm.io/gorm"
)

type ArticlesController struct {
}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// 获取参数
	id := route.GetRouteVariable("id", r)

	// 读取文章数据
	article, err := article.Get(id)

	// 错误处理
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 读取成功，显示文章
		// tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")
		// 增加删除按钮
		tmpl, err := template.New("show.gohtml").
			Funcs(template.FuncMap{
				"RouteName2URL": route.Name2URL,
				"Uint64ToString": types.Uint64ToString,
			}).ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)

		err = tmpl.Execute(w, article)
		logger.LogError(err)
	}
}
