package view

import (
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/wangyaodream/gerty-goblog/pkg/logger"
	"github.com/wangyaodream/gerty-goblog/pkg/route"
)

func Render(w io.Writer, name string, data interface{}) {
	// 设置模板相对路径
	viewDir := "resources/views/"

	// example.foo => example/foo
	name = strings.Replace(name, ".", "/", -1)

	files, err := filepath.Glob(viewDir + "layouts/*.gohtml")
	logger.LogError(err)

	newFiles := append(files, viewDir+name+".gohtml")

	tmpl, err := template.New(name + ".gohtml").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(newFiles...)
	logger.LogError(err)

	err = tmpl.ExecuteTemplate(w, "app", data)
	logger.LogError(err)
}
