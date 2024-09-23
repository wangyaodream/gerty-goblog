package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wangyaodream/gerty-goblog/app/http/controllers"
	"github.com/wangyaodream/gerty-goblog/app/http/middlewares"
)

func RegisterWebRoutes(r *mux.Router) {

	pc := new(controllers.PagesController)
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)

	// article
	ac := new(controllers.ArticlesController)
	r.HandleFunc("/", ac.Index).Methods("GET").Name("home")
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")

	r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")
	r.HandleFunc("/articles/create", middlewares.Auth(ac.Create)).Methods("GET").Name("articles.create")
	r.HandleFunc("/articles", middlewares.Auth(ac.Store)).Methods("POST").Name("articles.store")

	r.HandleFunc("/articles/{id:[0-9]+}/edit", middlewares.Auth(ac.Edit)).Methods("GET").Name("articles.edit")
	r.HandleFunc("/articles/{id:[0-9]+}", middlewares.Auth(ac.Update)).Methods("POST").Name("articles.update")

	r.HandleFunc("/articles/{id:[0-9]+}/delete", middlewares.Auth(ac.Delete)).Methods("POST").Name("articles.delete")

	// static route
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))

	// verify user
	auc := new(controllers.AuthController)
	r.HandleFunc("/auth/register", middlewares.Guest(auc.Register)).Methods("GET").Name("auth.register")
	r.HandleFunc("/auth/do-register", middlewares.Guest(auc.DoRegister)).Methods("POST").Name("auth.doregister")

	// login
	r.HandleFunc("/auth/login", middlewares.Guest(auc.Login)).Methods("GET").Name("auth.login")
	r.HandleFunc("/auth/dologin", middlewares.Guest(auc.DoLogin)).Methods("POST").Name("auth.dologin")

	// logout
	r.HandleFunc("/auth/logout", middlewares.Auth(auc.Logout)).Methods("POST").Name("auth.logout")

	// about user
	uc := new(controllers.UserController)
	r.HandleFunc("/users/{id:[0-9]+}", uc.Show).Methods("GET").Name("users.show")

	// 注册中间件
	r.Use(middlewares.StartSession)
	// use middleware
	// r.Use(middlewares.ForceHTML)

}
