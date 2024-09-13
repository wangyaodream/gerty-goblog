package main

import (
	"net/http"

	"github.com/wangyaodream/gerty-goblog/app/http/middlewares"
	"github.com/wangyaodream/gerty-goblog/bootstrap"
	"github.com/wangyaodream/gerty-goblog/pkg/logger"
)

func main() {
	bootstrap.SetupDB()
	router := bootstrap.SetupRoute()

	err := http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
	logger.LogError(err)
}
