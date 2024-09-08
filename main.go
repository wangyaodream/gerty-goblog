package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

// 包级别的变量不能使用:=表达式
// router := mux.NewRouter()
var router = mux.NewRouter()

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Hello there! Welcome!</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到:(</h1>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "article id: "+id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "The article list")
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "请提供正确的数据")
		return
	}

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := make(map[string]string)

	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度不正确(3-40)"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度不正确，不能低于10个字符"
	}

	if len(errors) == 0 {
		fmt.Fprint(w, "验证通过！<br>")
		fmt.Fprintf(w, "title: %s <br>", title)
		fmt.Fprintf(w, "title length: %d <br>", len(title))
		fmt.Fprintf(w, "body: %s <br>", body)
		fmt.Fprintf(w, "body length: %d <br>", len(body))
	} else {
		storeURL, _ := router.Get("articles.store").URL()

		data := ArticlesFormData{
			URL:    storeURL,
			Errors: errors,
			Title:  title,
			Body:   body,
		}
		t, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		err = t.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}

	// // r.Form比r.PostForm多了URL里参数的数据
	// fmt.Fprintf(w, "POST PostForm: %v <br>", r.PostForm)
	// fmt.Fprintf(w, "POST Form: %v <br>", r.Form)
	// fmt.Fprintf(w, "title : %v", title)
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

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{
		URL: storeURL,
	}
	t, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, data)
	if err != nil {
		panic(err)
	}
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

func main() {

	// router := http.NewServeMux()

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	// custom 404 page
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// use middleware
	router.Use(forceHTMLMiddleware)

	// use naming router
	homeURL, _ := router.Get("home").URL()
	fmt.Println("homeURL: ", homeURL)
	articleURL, _ := router.Get("articles.show").URL("id", "123")
	fmt.Println("articleURL: ", articleURL)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
