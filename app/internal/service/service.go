package service

import (
	"CommentTree/internal/entities/comment"
	"CommentTree/internal/storage"
	"errors"
	"fmt"
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

func (s *Service) CreateComment(c comment.Comment) (int64, error) {
	const op = "internal.service.Create"

	if c.Message == "" {
		return 0, fmt.Errorf("%w: %s", ErrWrongData, "empty comment text")
	}
	if c.ParentID < 0 {
		return 0, fmt.Errorf("%w: %s", ErrWrongData, "wrong parent id")
	}

	id, err := s.str.CreateComment(c)
	if errors.Is(err, storage.ErrWrongForeignKey) {
		return 0, fmt.Errorf("%w: %s", ErrWrongData, "wrong parent id")
	} else if err != nil {
		return 0, fmt.Errorf("%s: (%w)%w", op, ErrStorageInternal, err)
	}

	return id, nil
}

func (s *Service) Comments(parentID int64, opts *comment.GetterOpts) (*comment.CommentView, error) {
	const op = "internal.service.comments"

	// if parentID == 0, select * without conditions
	if parentID < 0 {
		return nil, fmt.Errorf("%w: %s", ErrWrongData, "wrong id")
	}

	result, err := s.str.Comments(parentID, opts)
	if errors.Is(err, storage.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("%s: (%w)%w", op, ErrStorageInternal, err)
	}

	return result, nil
}

func (s *Service) DeleteComment(id int64) error {
	const op = "internal.service.Delete"

	if id <= 0 {
		return fmt.Errorf("%w: %s", ErrWrongData, "wrong id")
	}

	err := s.str.DeleteComment(id)
	if errors.Is(err, storage.ErrNotAffected) {
		return ErrNotAffected
	} else if err != nil {
		return fmt.Errorf("%s: (%w)%w", op, ErrStorageInternal, err)
	}

	return nil
}
