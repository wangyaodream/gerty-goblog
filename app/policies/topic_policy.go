package policies

import (
	"github.com/wangyaodream/gerty-goblog/app/models/article"
	"github.com/wangyaodream/gerty-goblog/pkg/auth"
)

// 定义授权策略
func CanModifyArticle(_article article.Article) bool {
	return auth.User().ID == _article.User.ID
}
