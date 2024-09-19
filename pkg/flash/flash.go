package flash

import (
	"encoding/gob"

	"github.com/wangyaodream/gerty-goblog/pkg/session"
)

type Flashes map[string]interface{}

var flashkey = "_flashes"

func init() {
	gob.Register(Flashes{})
}

func Info(message string) {
	addFlash("info", message)
}

func Warning(message string) {
	addFlash("warning", message)
}

func Success(message string) {
	addFlash("success", message)
}

func Danger(message string) {
	addFlash("danger", message)
}

// All returns all flash messages
func All() Flashes {
	val := session.Get(flashkey)
	flashMessages, ok := val.(Flashes)
	if !ok {
		return nil
	}
	session.Forget(flashkey)
	return flashMessages
}

func addFlash(key string, message string) {
	flashes := Flashes{}
	flashes[key] = message
	session.Put(flashkey, flashes)
	session.Save()
}
