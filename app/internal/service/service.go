package service

import (
	"errors"

	"CommentTree/internal/entities/comment"
)

var (
	ErrWrongData       = errors.New("wrong data")
	ErrNotAffected     = errors.New("not affected")
	ErrNotFound        = errors.New("not found")
	ErrStorageInternal = errors.New("internal storage error")
)

type str interface {
	CreateComment(c comment.Comment) (int64, error)
	Comments(parentID int64, opts *comment.GetterOpts) (*comment.CommentView, error)
	DeleteComment(id int64) error
}

type Service struct {
	str
}

func New(str str) *Service {
	return &Service{
		str: str,
	}
}

func (s *Service) Shutdown() {}
