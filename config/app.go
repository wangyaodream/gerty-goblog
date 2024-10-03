package config

import (
	"github.com/wangyaodream/gerty-goblog/pkg/config"
)

func init() {
	config.Add("app", config.StrMap{
		"name":   config.Env("APP_NAME", "GoBlog"),
		"env":    config.Env("APP_ENV", "production"),
		"debug":  config.Env("APP_DEBUG", false),
		"port":   config.Env("APP_PORT", "3000"),
		"key":    config.Env("APP_KEY", "33446a9dcf9ea060a0a6532b166da32f304af0de"),
		"url":    config.Env("APP_URL", "http://localhost3000"),
		"apikey": config.Env("BIGMODEL_APIKEY", nil),
	})
}
