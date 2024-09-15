package view

import (
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/wangyaodream/gerty-goblog/pkg/logger"
	"github.com/wangyaodream/gerty-goblog/pkg/route"
)

type D map[string]interface{}

func RenderTemplate(w io.Writer, name string, data interface{}, tplFiles ...string) {
	// 设置模板相对路径
	viewDir := "resources/views/"

	// 支持多模版
	for i, f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}

	layoutFiles, err := filepath.Glob(viewDir + "layouts/*.gohtml")
	logger.LogError(err)

	// 合并所有模版文件
	allFiles := append(layoutFiles, tplFiles...)

	// 解析所有模版文件
	tmpl, err := template.New("").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(allFiles...)
	logger.LogError(err)

	err = tmpl.ExecuteTemplate(w, name, data)
	logger.LogError(err)
}

func RenderSimple(w io.Writer, data interface{}, tplFiles ...string) {
	RenderTemplate(w, "simple", data, tplFiles...)
}

func Render(w io.Writer, data interface{}, tplFiles ...string) {
	RenderTemplate(w, "app", data, tplFiles...)
}
