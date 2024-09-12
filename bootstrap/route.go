package bootstrap

import (
	"github.com/gorilla/mux"
	"github.com/wangyaodream/gerty-goblog/routes"
)

func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)
	return router
}
