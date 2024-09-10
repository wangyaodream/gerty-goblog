package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// 包级别的变量不能使用:=表达式
// router := mux.NewRouter()
var router = mux.NewRouter()
var db *sql.DB

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func (a Article) Link() string {
    showURL, err := router.Get("articles.show").URL("id", strconv.FormatInt(a.ID, 10))
    if err != nil {
        checkError(err)
        return ""
    }
    return showURL.String()
}

func initDB() {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	dbType := os.Getenv("DB_TYPE")
	mysqlConnStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/goblog",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"))

	// PostgreSQL连接字符串
	postgresConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	switch dbType {
	case "mysql":
		db, err = connectDB("mysql", mysqlConnStr)
		checkError(err)

		// 配置连接属性
		db.SetMaxOpenConns(25)                 // 最大连接数
		db.SetMaxIdleConns(25)                 // 最大空闲数
		db.SetConnMaxLifetime(5 * time.Minute) // 每个链接的过期时间
	case "postgresql":
		db, err = connectDB("postgres", postgresConnStr)
		checkError(err)
	}
}

func connectDB(dbType, connStr string) (*sql.DB, error) {
	db, err := sql.Open(dbType, connStr)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
); `

	_, err := db.Exec(createArticlesSQL)
	checkError(err)
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

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT * from articles")
    checkError(err)
    defer rows.Close()

    var articles []Article

    for rows.Next() {
        var article Article

        err := rows.Scan(&article.ID, &article.Title, &article.Body)
        checkError(err)

        articles = append(articles, article)
    }

    err = rows.Err()
    checkError(err)

    // 加载模版
    tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
    checkError(err)

    // 渲染模版，将articles中的数据传入到模版中
    err = tmpl.Execute(w, articles)
    checkError(err)
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "请提供正确的数据")
		return
	}

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

    errors := validateArticleFormData(title, body)

	if len(errors) == 0 {
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprint(w, "插入成功，ID为"+strconv.FormatInt(lastInsertID, 10))
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}

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

func saveArticleToDB(title string, body string) (int64, error) {
	// 初始化变量
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)

	// 获取一个prepare声明
	stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?, ?)")
	if err != nil {
		return 0, err
	}

	// 在defer中声明关闭
	defer stmt.Close()

	// 执行请求
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return 0, err
	}

	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}

	return 0, err
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

// 对应数据库中读取的数据
type Article struct {
	Title, Body string
	ID          int64
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	// 获取参数
    id := getRouteVariable("id", r)

	// 读取文章数据
    article, err := getArticleByID(id)

	// 错误处理
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
        // 读取成功，显示文章
        tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")
        checkError(err)

        err = tmpl.Execute(w, article)
        checkError(err)
	}
}

func getRouteVariable(parameterName string, r *http.Request) string {
    vars := mux.Vars(r)
    return vars[parameterName]
}

func getArticleByID(id string) (Article, error) {
    article := Article{}
    query := "SELECT * FROM articles WHERE id = ?"
    err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
    return article, err
}

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
    // Get URLs param
    id := getRouteVariable("id", r)

    // read the data of post
    article, err := getArticleByID(id)

    // handle error
    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        } else {
            checkError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        }
    } else {
        // 读取成功
        updateURL, _ := router.Get("articles.update").URL("id", id)
        data := ArticlesFormData{
            Title: article.Title,
            Body: article.Body,
            URL: updateURL,
            Errors: nil,
        }
        tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
        checkError(err)

        err = tmpl.Execute(w, data)
        checkError(err)
    }

}

func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
    id := getRouteVariable("id", r)
    // 检查该文章
    _, err := getArticleByID(id)

    if err != nil {
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusNotFound)
            fmt.Fprint(w, "404 文章未找到")
        } else {
            checkError(err)
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprint(w, "500 服务器内部错误")
        }
    } else {
        // 没有错误的情况
        title := r.PostFormValue("title")
        body := r.PostFormValue("body")

        errors := validateArticleFormData(title, body)

        // check errors
        if len(errors) == 0 {
            query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"
            res, err := db.Exec(query, title, body, id)

            if err != nil {
                checkError(err)
                w.WriteHeader(http.StatusInternalServerError)
                fmt.Fprint(w, "500 服务器内部错误")
            }

            // 更新成功，跳转到文章详情页
            if n, _ := res.RowsAffected(); n > 0 {
                showURL, _ := router.Get("articles.show").URL("id", id)
                http.Redirect(w, r, showURL.String(), http.StatusFound)
            } else {
                fmt.Fprint(w, "没有任何修改!")
            }
            
        } else {
            // 有错误，将错误输出
            updateURL, _ := router.Get("articles.update").URL("id", id)
            data := ArticlesFormData{
                Title: title,
                Body: body,
                URL: updateURL,
                Errors: errors,
            }
            tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
            checkError(err)

            err = tmpl.Execute(w, data)
            checkError(err)
        }

        
    }
}

func main() {
	// 初始化数据库
	initDB()
	createTables()
	defer db.Close()

	// router := http.NewServeMux()

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")
    router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
    router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")

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
