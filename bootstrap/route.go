package bootstrap

import (
	"github.com/gorilla/mux"
	"github.com/wangyaodream/gerty-goblog/pkg/route"
	"github.com/wangyaodream/gerty-goblog/routes"
)

func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)

	route.SetRoute(router)

	return router
}
