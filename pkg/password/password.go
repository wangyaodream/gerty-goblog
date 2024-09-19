package password

import (
	"github.com/wangyaodream/gerty-goblog/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) string {
	// 第二个参数是cost，数值越大越费时
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	logger.LogError(err)

	return string(bytes)
}

func CheckHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	logger.LogError(err)
	return err == nil
}

func IsHashed(str string) bool {
	return len(str) == 60
}
