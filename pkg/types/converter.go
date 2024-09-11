package types

import "strconv"


func Int64ToString(num int64) string {
    // 将一个整型值转换成一个字符串，以十进制
    return strconv.FormatInt(num, 10)
}
