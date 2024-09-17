package requests

import (
	"github.com/thedevsaddam/govalidator"
	"github.com/wangyaodream/gerty-goblog/app/models/user"
)

func ValidateRegistrationForm(data user.User) map[string][]string {

	rules := govalidator.MapData{
		// not_exists是在request.go中定义的，传入的参数是表名和字段名
		"name":             []string{"required", "alpha_num", "between:3,20", "not_exists:users,name"},
		"email":            []string{"required", "min:4", "max:30", "email", "not_exists:users,email"},
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
	}

	// 自定义错误消息
	messages := govalidator.MapData{
		"name": []string{
			"required:用户名不能为空",
			"alpha_num:用户名格式不正确, 只能是字母数字",
			"between:用户名长度必须在3~20之间",
		},
		"email": []string{
			"required:邮箱不能为空",
			"min:邮箱长度不能小于4",
			"max:邮箱长度不能大于30",
		},
		"password": []string{
			"required:密码不能为空",
			"min:密码长度不能小于6",
		},
		"password_confirm": []string{
			"required:确认密码不能为空",
		},
	}

	// 配置初始化
	options := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: "valid", // 默认是 validation
	}

	// 验证器
	errs := govalidator.New(options).ValidateStruct()

	if data.Password != data.PasswordConfirm {
		errs["password_confirm"] = append(errs["password_confirm"], "两次输入密码不一致")
	}

	return errs
}
