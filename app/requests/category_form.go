package requests

import (
	"github.com/thedevsaddam/govalidator"
	"github.com/wangyaodream/gerty-goblog/app/models/category"
)

func ValidateCategoryForm(data category.Category) map[string][]string {
	rules := govalidator.MapData{
		"name": []string{"required", "min:2", "max:8", "not_exists:categories,name"},
	}

	messages := govalidator.MapData{
		"name": []string{
			"required:分类名称为必填项",
			"min:分类名称长度需大于2",
			"max:分类名称长度需小于8",
		},
	}

	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      messages,
	}

	return govalidator.New(opts).ValidateStruct()
}
