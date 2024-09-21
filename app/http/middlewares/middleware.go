package middlewares

import "net/http"

type HTTPHandlerFunc func(http.ResponseWriter, *http.Request)
