package middlewares

import "net/http"

func ForceHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置响应头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 将请求传递下去
		next.ServeHTTP(w, r)
	})
}
