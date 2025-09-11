package storage

import (
	"CommentTree/internal/entities/comment"
	"CommentTree/pkg/errs"
	"errors"
	"fmt"
)

var (
	ErrNotAffected = errors.New("not affected")
)

type db interface {
	Create(c comment.Comment) (int64, error)
	Comments(id int64, opts *comment.GetterOpts) ([]comment.Comment, error)
	Delete(id int64) error
}

type Storage struct {
	db
}

func New(db db) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Create(c comment.Comment) (int64, error) {
	const op = "internal.storage.Create"

	id, err := s.db.Create(c)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) Comments(id int64, opts *comment.GetterOpts) ([]comment.Comment, error) {
	const op = "internal.storage.Comments"
	
	return nil, nil
}

func (s *Storage) Delete(id int64) error {
	const op = "internal.storage.Delete"

	err := s.db.Delete(id)
	if errors.Is(err, errs.ErrDBNotAffected) {
		return ErrNotAffected
	} else if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
