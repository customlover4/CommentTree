package main

import (
	"CommentTree/internal/service"
	"CommentTree/internal/storage"
	"CommentTree/internal/storage/postgres"
	"CommentTree/internal/web"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

var (
	ConnString = "postgresql://dev:qqq@localhost:5432/test?sslmode=disable"
	ConfigPath = "../config/config.yml"
	Port       = "8080"

	// TODO: move it to config.
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 5 * time.Second
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

	tmp = os.Getenv("WEB_PORT")
	if tmp != "" {
		Port = tmp
	}
}

func shutdown(server *http.Server, srv *service.Service, str *storage.Storage) {
	_ = server.Shutdown(context.Background())
	srv.Shutdown()
	str.Shutdown()
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	cfg := config.New()
	err := cfg.Load(ConfigPath)
	if err != nil {
		panic(err)
	}

	if cfg.GetString("debug") == "false" {
		gin.SetMode(gin.ReleaseMode)
	}

	p := postgres.New(ConnString)
	str := storage.New(p)
	srv := service.New(str)

	router := ginext.New()
	web.Routes(router, srv)
	server := &http.Server{
		Handler:      router,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
	}

	l := net.ListenConfig{}

	ctx, finish := context.WithCancel(context.Background())
	defer finish()

	listener, err := l.Listen(ctx, "tcp", ":"+Port)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := server.Serve(listener); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				fmt.Fprintln(os.Stderr, err)
				sig <- syscall.SIGINT
			}
		}
	}()

	zlog.Logger.Info().Msg("server started on :" + Port)

	<-sig

	finish()
	shutdown(server, srv, str)

	zlog.Logger.Info().Msg("server stopped")
}
