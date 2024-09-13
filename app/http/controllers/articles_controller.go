package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"unicode/utf8"

	"github.com/wangyaodream/gerty-goblog/app/models/article"
	"github.com/wangyaodream/gerty-goblog/pkg/logger"
	"github.com/wangyaodream/gerty-goblog/pkg/route"
	"github.com/wangyaodream/gerty-goblog/pkg/types"
	"gorm.io/gorm"
)

type ArticlesController struct {
}

type ArticlesFormData struct {
    Title, Body string
    URL     string
    Errors map[string]string
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

func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
    
    // 交给app/models/article中的crud处理
    articles, err := article.GetAll()

    if err != nil {
        // 这里错误是呈现在后台
        logger.LogError(err)
        w.WriteHeader(http.StatusInternalServerError)
        // 这里的信息呈现在网页
        fmt.Fprint(w, "500 服务器内部错误")
    } else {
        tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
        logger.LogError(err)

        // render template
        err = tmpl.Execute(w, articles)
        logger.LogError(err)
    }
}

func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {

    storeURL := route.Name2URL("articles.store")
    data := ArticlesFormData{
        Title: "",
        Body: "",
        URL: storeURL,
        Errors: nil,
    }
    tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
    if err != nil {
        panic(err)
    }

    // render create page
    err = tmpl.Execute(w, data)
    if err != nil {
        panic(err)
    }
}

func validateArticleFormData(title string, body string) map[string]string {
    errors := make(map[string]string)

    // valid title
    if title == "" {
        errors["title"] = "标题不能为空"
    } else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
        errors["title"] = "标题长度需要介于 3-40"
    }

    // valid body
    if body == "" {
        errors["body"] = "内容不能为空"
    } else if utf8.RuneCountInString(body) < 10 {
        errors["body"] = "内容长度不能低于10"
    }
    return errors
}

func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
    title := r.PostFormValue("title")
    body := r.PostFormValue("body")

    errors := validateArticleFormData(title, body)
    
    if len(errors) == 0 {
        _article := article.Article{
            Title: title,
            Body: body,
        }
        _article.Create()
        if _article.ID > 0 {
            fmt.Fprint(w, "插入成功! ID: " + strconv.FormatUint(_article.ID, 10))
        } else {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "创建文章失败!")
        }
    } else {

        storeURL := route.Name2URL("articles.store")

        data := ArticlesFormData{
            Title: title,
            Body: body,
            URL: storeURL,
            Errors: errors,
        }

        tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
        logger.LogError(err)
        // render create page with invalid message
        err = tmpl.Execute(w, data)
        logger.LogError(err)
    }

}

// Edit
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
    // get parameter by 
    id := route.GetRouteVariable("id", r)

    article, err := article.Get(id)

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
        updateURL := route.Name2URL("articles.update", "id", id)
        data := ArticlesFormData{
            Title: article.Title,
            Body: article.Body,
            URL: updateURL,
            Errors: nil,
        }
        tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
        logger.LogError(err)

        err = tmpl.Execute(w, data)
        logger.LogError(err)
    }
}

func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
    id := route.GetRouteVariable("id", r)

    _article, err := article.Get(id)

    if err != nil {
        if err == gorm.ErrRecordNotFound {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        } else {
            logger.LogError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错")
        }
    } else {
        title := r.PostFormValue("title")
        body := r.PostFormValue("body")

        errors := validateArticleFormData(title, body)

        if len(errors) == 0 {
            _article.Title = title
            _article.Body = body

            rowsAffected, err := _article.Update()

            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprint(w, "500 服务器内部错误")
            }

            if rowsAffected > 0 {
                showURL := route.Name2URL("articles.show", "id", id)
                http.Redirect(w, r, showURL, http.StatusFound)
            } else {
                fmt.Fprint(w, "没做任何修改！")
            }
        } else {

            updateURL := route.Name2URL("articles.update", "id", id)

            data := ArticlesFormData{
                Title: title,
                Body: body,
                URL: updateURL,
                Errors: errors,
            }
            tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
            logger.LogError(err)

            err = tmpl.Execute(w, data)
            logger.LogError(err)

        }
    }
}
