package user

import (
	"github.com/wangyaodream/gerty-goblog/app/models"
	"github.com/wangyaodream/gerty-goblog/pkg/route"
)

type User struct {
	models.BaseModel

	Name     string `gorm:"type:varchar(255);not null;unique" valid:"name"`
	Email    string `gorm:"type:varchar(255);unique" valid:"email"`
	Password string `gorm:"type:varchar(255)" valid:"password"`
	// 设置"-"表示GORM读写时略过字段，仅用于验证
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}

func (u User) Link() string {
	return route.Name2URL("users.show", "id", u.GetStringID())
}
