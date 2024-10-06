package logger

import "log"

func LogError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func LogInfo(info string) {
	log.Println(info)
}
