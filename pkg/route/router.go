package route

import "github.com/gorilla/mux"


var Router *mux.Router

func Initialize() {
    Router = mux.NewRouter()
}


func Name2URL(routeName string, pairs ...string) string {
    url, err := Router.Get(routeName).URL(pairs...)
    if err != nil {
        return ""
    }

    return url.String()
}
