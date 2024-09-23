package view

import (
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/wangyaodream/gerty-goblog/app/models/user"
	"github.com/wangyaodream/gerty-goblog/pkg/auth"
	"github.com/wangyaodream/gerty-goblog/pkg/flash"
	"github.com/wangyaodream/gerty-goblog/pkg/logger"
	"github.com/wangyaodream/gerty-goblog/pkg/route"
)

type D map[string]interface{}

func RenderTemplate(w io.Writer, name string, data D, tplFiles ...string) {

	// 通用模板数据
	data["isLogined"] = auth.Check()
	data["loginUser"] = auth.User
	data["flash"] = flash.All()
	data["Users"], _ = user.All()

	// 合并所有模版文件
	allFiles := getTemplateFiles(tplFiles...)

	// 解析所有模版文件
	tmpl, err := template.New("").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(allFiles...)
	logger.LogError(err)

	err = tmpl.ExecuteTemplate(w, name, data)
	logger.LogError(err)
}

func RenderSimple(w io.Writer, data D, tplFiles ...string) {
	RenderTemplate(w, "simple", data, tplFiles...)
}

func Render(w io.Writer, data D, tplFiles ...string) {
	RenderTemplate(w, "app", data, tplFiles...)
}

func getTemplateFiles(tplFiles ...string) []string {
	viewDir := "resources/views/"

	for i, f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}

	layoutFiles, err := filepath.Glob(viewDir + "layouts/*.gohtml")
	logger.LogError(err)

	return append(layoutFiles, tplFiles...)
}
