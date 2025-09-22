package handlers

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

const (
	MainPageHTML      = "index.html"
	CreateCommentHTML = "create.html"
	ShowCommentsHTML  = "show.html"
	CommentsHTML      = "find.html"
)

func MainPage(ctx *ginext.Context) {
	ctx.HTML(http.StatusOK, MainPageHTML, nil)
}

func CreateCommentPage(ctx *ginext.Context) {
	ctx.HTML(http.StatusOK, CreateCommentHTML, nil)
}

func ShowCommentsPage(ctx *ginext.Context) {
	ctx.HTML(http.StatusOK, ShowCommentsHTML, nil)
}

func CommentsPage(ctx *ginext.Context) {
	ctx.HTML(http.StatusOK, CommentsHTML, nil)
}
