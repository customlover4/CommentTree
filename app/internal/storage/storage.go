package storage

import (
	"errors"

	"CommentTree/internal/entities/comment"
)

var (
	ErrNotAffected     = errors.New("not affected")
	ErrNotFound        = errors.New("not found")
	ErrWrongForeignKey = errors.New("wrong foreign key")
)

type db interface {
	CreateComment(c comment.Comment) (int64, error)
	Parent(id int64) (comment.Comment, error)
	Childs(parentID int64, opts *comment.GetterOpts) ([]comment.Comment, error)
	DeleteComment(id int64) error

	Shutdown()
}

type Storage struct {
	db
}

func New(db db) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Shutdown() {
	s.db.Shutdown()
}
