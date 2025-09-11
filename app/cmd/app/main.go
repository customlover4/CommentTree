package main

import (
	"CommentTree/internal/entities/comment"
	"CommentTree/internal/service"
	"CommentTree/internal/storage"
	"CommentTree/internal/storage/postgres"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/zlog"
)

var (
	ConnString = "postgresql://dev:qqq@localhost:5432/test?sslmode=disable"
	ConfigPath = "../config/config.yml"
)

func init() {
	zlog.Init()

	tmp := os.Getenv("POSTGRES_CONN")
	if tmp != "" {
		ConnString = tmp
	}

	tmp = os.Getenv("CONFIG_PATH")
	if tmp != "" {
		ConfigPath = tmp
	}
}

func main() {
	cfg := config.New()
	err := cfg.Load(ConfigPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if cfg.GetString("debug") == "false" {
		gin.SetMode(gin.ReleaseMode)
	}

	p := postgres.New(ConnString)
	str := storage.New(p)
	// p.Create(comment.Comment{Message: "привет бро"})
	// p.Create(comment.Comment{Message: "здарова", ParentID: 4})
	// p.Create(comment.Comment{Message: "здарова опять", ParentID: 4})
	fmt.Println(p.Comments(0, &comment.GetterOpts{
		Substr: "здарова",
	}))
	srv := service.New(str)

	_ = srv

}
