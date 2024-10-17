package controllers

import (
	"net/http"

	"github.com/wangyaodream/gerty-goblog/app/models/comment"
	"github.com/wangyaodream/gerty-goblog/pkg/logger"
	"github.com/wangyaodream/gerty-goblog/pkg/route"
	"github.com/wangyaodream/gerty-goblog/pkg/view"
)

type CommentsController struct {
}

func (*CommentsController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	comments, err := comment.Get(id)
	if err != nil {
		logger.LogError(err)
	}

	view.Render(w, view.D{
		"Comment": comments,
	}, "comments.show")

}
