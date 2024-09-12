package controllers

import (
	"fmt"
	"net/http"
)

type PagesController struct {
}

func (*PagesController) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello Welcome to my Blog!</h1>")
}

func (*PagesController) About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>About Us</h1>")
}

// NotFound page
func (*PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>Page Not Found :(</h1>")
}
