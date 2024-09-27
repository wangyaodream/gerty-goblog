package requests

import (
	"github.com/thedevsaddam/govalidator"
	"github.com/wangyaodream/gerty-goblog/app/models/article"
)

func ValidateArticleForm(data article.Article) map[string][]string {
	rules := govalidator.MapData{
		"title": []string{"required", "min_cn:3", "max_cn:40"},
		"body":  []string{"required", "min_cn:10"},
	}

	// 定制错误消息
	messages := govalidator.MapData{
		"title": []string{
			"required:标题为必填项",
			"min:标题长度需大于3",
			"max:标题长度需小于40",
		},
		"body": []string{
			"required:文章内容不能为空",
			"min:文章内容长度需大于10",
		},
	}

	// 配置初始化
	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      messages,
	}

	return govalidator.New(opts).ValidateStruct()
}
