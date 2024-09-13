package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gorilla/mux"
	"github.com/wangyaodream/gerty-goblog/bootstrap"
	"github.com/wangyaodream/gerty-goblog/pkg/database"
	"github.com/wangyaodream/gerty-goblog/pkg/logger"
)

// 包级别的变量不能使用:=表达式
// router := mux.NewRouter()
var router *mux.Router
var DB *sql.DB = database.DB

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}


func (a Article) Delete() (rowsAffected int64, err error) {
	res, err := DB.Exec("DELETE FROM articles WHERE id = ?", strconv.FormatInt(a.ID, 10))

	if err != nil {
		return 0, err
	}

	// 删除成功
	if n, _ := res.RowsAffected(); n > 0 {
		return n, nil
	}

	return 0, nil
}

func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)

	// valid title
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需要介于3-40"
	}

	// valid body
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需要大于等于10"
	}
	return errors
}


// force HTML content type
func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set the content type
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// call the next handler
		next.ServeHTTP(w, r)
	})
}

// 处理URL最后的斜杠
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 排除首页的斜杠
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		// 将请求传递下去
		next.ServeHTTP(w, r)
	})
}

// 对应数据库中读取的数据
type Article struct {
	Title, Body string
	ID          int64
}

func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := DB.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}


func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)
	article, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		rowsAffected, err := article.Delete()

		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			// 未发生错误
			if rowsAffected > 0 {
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
}

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

func main() {
	// 初始化数据库
	database.Initialize()
	DB = database.DB


    bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	// router := http.NewServeMux()

	// router remove
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

	// use middleware
	router.Use(forceHTMLMiddleware)

	// use naming router
	homeURL, _ := router.Get("home").URL()
	fmt.Println("homeURL: ", homeURL)
	articleURL, _ := router.Get("articles.show").URL("id", "123")
	fmt.Println("articleURL: ", articleURL)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
