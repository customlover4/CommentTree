package web

import (
	"CommentTree/internal/service"
	"CommentTree/internal/web/handlers"

	"github.com/wb-go/wbf/ginext"
)

func Routes(router *ginext.Engine, s *service.Service) {

	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", handlers.MainPage())
	router.POST("/comments", handlers.CreateComment(s))
	router.DELETE("/comments/:id", handlers.DeleteComment(s))
}
