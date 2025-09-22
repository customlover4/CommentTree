package web

import (
	"CommentTree/internal/service"
	"CommentTree/internal/web/handlers"

	"github.com/wb-go/wbf/ginext"
)

func Routes(router *ginext.Engine, s *service.Service) {
	router.Delims("__", "__").LoadHTMLGlob("templates/*.html")

	// html
	router.GET("/", handlers.MainPage)
	router.GET("/create", handlers.CreateCommentPage)
	// router.GET("/show", handlers.ShowCommentsPage)
	router.GET("/show", handlers.CommentsPage)

	// api
	router.POST("/comments", handlers.CreateComment(s))
	router.DELETE("/comments/:id", handlers.DeleteComment(s))
	router.GET("/comments", handlers.Comments(s))
}
