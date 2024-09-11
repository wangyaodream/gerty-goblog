package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wangyaodream/gerty-goblog/pkg/logger"
)


var Router *mux.Router

func Initialize() {
    Router = mux.NewRouter()
}


func Name2URL(routeName string, pairs ...string) string {
    url, err := Router.Get(routeName).URL(pairs...)
    if err != nil {
        logger.LogError(err)
        return ""
    }

    return url.String()
}

func GetRouteVariable(parameterName string, r *http.Request) string {
    vars := mux.Vars(r)
    return vars[parameterName]
}
