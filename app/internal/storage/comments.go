package storage

import "CommentTree/internal/entities/comment"

type DB interface {
	Create(c comment.Comment) (int64, error)
	Comments(id int64) ([]comment.Comment, error)
	Delete(id int64) error
}

type CommentStorage struct {
	db DB
}

func NewComment(db DB) *CommentStorage {
	return &CommentStorage{
		db: db,
	}
}

func (cs *CommentStorage) Create(c comment.Comment) (int64, error) {
	return 0, nil
}

func (cs *CommentStorage) Comments(id int64) ([]comment.Comment, error) {
	return nil, nil
}

func (cs *CommentStorage) Delete(id int64) error {
	return nil
}
