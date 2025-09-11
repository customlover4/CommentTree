package storage

import "CommentTree/internal/entities/comment"

type db interface {
	Create(c comment.Comment) (int64, error)
	Comments(id int64) ([]comment.Comment, error)
	Delete(id int64) error
}

type Storage struct {
	db
}

func NewComments(db db) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Create(c comment.Comment) (int64, error) {
	return 0, nil
}

func (s *Storage) Comments(id int64) ([]comment.Comment, error) {
	return nil, nil
}

func (s *Storage) Delete(id int64) error {
	return nil
}
