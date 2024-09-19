package session

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/wangyaodream/gerty-goblog/pkg/logger"
)

var Store = sessions.NewCookieStore([]byte("33446a9dcf9ea060a0a6532b166da32f304af0de"))

var Session *sessions.Session
var Request *http.Request
var Response http.ResponseWriter

func StartSession(w http.ResponseWriter, r *http.Request) {
	var err error

	Session, err = Store.Get(r, "goblog-session")
	logger.LogError(err)

	Request = r
	Response = w
}

func Put(key string, value interface{}) {
	Session.Values[key] = value
	Save()
}

func Get(key string) interface{} {
	return Session.Values[key]
}

func Forget(key string) {
	delete(Session.Values, key)
	Save()
}

func Flush() {
	Session.Options.MaxAge = -1
	Save()
}

func Save() {
	err := Session.Save(Request, Response)
	logger.LogError(err)
}
