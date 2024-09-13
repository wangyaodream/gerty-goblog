package types

import (
	"strconv"

	"github.com/wangyaodream/gerty-goblog/pkg/logger"
)


func Int64ToString(num int64) string {
    // 将一个整型值转换成一个字符串，以十进制
    return strconv.FormatInt(num, 10)
}

func Uint64ToString(num uint64) string {
    return strconv.FormatUint(num, 10)
}

func StringToUint64(str string) uint64 {
    i, err := strconv.ParseUint(str, 10, 64)
    if err != nil {
        logger.LogError(err)
    }
    return i
}
