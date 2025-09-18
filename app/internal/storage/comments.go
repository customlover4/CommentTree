package storage

import (
	"CommentTree/internal/entities/comment"
	"CommentTree/pkg/errs"
	"errors"
	"fmt"
)

func (s *Storage) CreateComment(c comment.Comment) (int64, error) {
	const op = "internal.storage.Create"

	id, err := s.db.CreateComment(c)
	err = s.UnwrapError(err)
	if errors.Is(err, errs.ErrDBViolatesForeignKey) {
		return 0, ErrWrongForeignKey
	} else if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) Comments(parentID int64, opts *comment.GetterOpts) (*comment.CommentView, error) {
	const op = "internal.storage.Comments"

	var parent comment.Comment
	var err error
	if parentID != 0 {
		parent, err = s.Parent(parentID)
		if errors.Is(err, errs.ErrDBNotFound) {
			return nil, ErrNotFound
		} else if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	childs, err := s.Childs(parentID, opts)
	if errors.Is(err, errs.ErrDBNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &comment.CommentView{
		ParentComment: parent,
		Childs:        childs,
	}, nil
}

func (s *Storage) DeleteComment(id int64) error {
	const op = "internal.storage.Delete"

	err := s.db.DeleteComment(id)
	if errors.Is(err, errs.ErrDBNotAffected) {
		return ErrNotAffected
	} else if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
