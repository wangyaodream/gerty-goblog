package view

import (
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/wangyaodream/gerty-goblog/pkg/auth"
	"github.com/wangyaodream/gerty-goblog/pkg/logger"
	"github.com/wangyaodream/gerty-goblog/pkg/route"
)

type D map[string]interface{}

func RenderTemplate(w io.Writer, name string, data D, tplFiles ...string) {

	// 判断登录,直接查看session中的uid值是否存在
	data["isLogined"] = auth.Check()

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
