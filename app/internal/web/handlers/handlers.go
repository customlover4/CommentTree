package handlers

import (
	"CommentTree/internal/entities/comment"
	"CommentTree/internal/entities/request"
	"CommentTree/internal/entities/response"
	"CommentTree/internal/service"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

const (
	MainPageHTML = "index.html"

	InternalError = "internal error on service"
)

type servicer interface {
	CreateComment(c comment.Comment) (int64, error)
	Comments(parentID int64, opts *comment.GetterOpts) (*comment.CommentView, error)
	DeleteComment(id int64) error
}

func MainPage() ginext.HandlerFunc {
	return func(ctx *ginext.Context) {
		ctx.HTML(http.StatusOK, MainPageHTML, nil)
	}
}

func CreateComment(s servicer) ginext.HandlerFunc {
	return func(ctx *ginext.Context) {
		const op = "internal.web.handlers.CreateComment"

		ctx.Header("Content-Type", "application/json")

		var req request.CreateComment
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, response.Error(
				"wrong json, data or types in json",
			))
			return
		}

		ok := req.Validate()
		if ok != "" {
			ctx.JSON(http.StatusBadRequest, response.Error(
				ok,
			))
			return
		}

		id, err := s.CreateComment(comment.Comment{
			Message:  req.Message,
			ParentID: req.ParentID,
		})
		if errors.Is(err, service.ErrWrongData) {
			ctx.JSON(http.StatusServiceUnavailable, response.Error(
				err.Error(),
			))
			return
		} else if err != nil {
			zlog.Logger.Error().Err(fmt.Errorf("%s: %w", op, err)).Send()
			ctx.JSON(http.StatusInternalServerError, response.Error(
				InternalError,
			))
			return
		}

		ctx.JSON(http.StatusOK, response.Result(id))
	}
}

func DeleteComment(s servicer) ginext.HandlerFunc {
	return func(ctx *ginext.Context) {
		const op = "internal.web.handlers.DeleteComment"

		ctx.Header("Content-Type", "application/json")

		idParam := ctx.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Error(
				"id should be numeric",
			))
			return
		}
		if id <= 0 {
			ctx.JSON(http.StatusBadRequest, response.Error(
				"id shoud be > 0",
			))
			return
		}

		err = s.DeleteComment(id)
		if errors.Is(err, service.ErrWrongData) {
			ctx.JSON(http.StatusServiceUnavailable, response.Error(
				err.Error(),
			))
		} else if errors.Is(err, service.ErrNotAffected) {
			ctx.JSON(http.StatusServiceUnavailable, response.Error(
				"can't find comment with this id",
			))
			return
		} else if err != nil {
			zlog.Logger.Error().Err(fmt.Errorf("%s: %w", op, err)).Send()
			ctx.JSON(http.StatusInternalServerError, response.Error(
				InternalError,
			))
			return
		}

		ctx.JSON(http.StatusOK, response.OK())
	}
}
