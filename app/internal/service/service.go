package service

import "CommentTree/internal/entities/comment"

type storage interface {
	Create(c comment.Comment) (int64, error)
	Comments(id int64) ([]comment.Comment, error)
	Delete(id int64) error
}

type Service struct {
	str storage
}

func New(str storage) *Service {
	return &Service{
		str: str,
	}
}

func (s *Service) Create(c comment.Comment) (int64, error) {
	return 0, nil
}

func (s *Service) Comments(id int64) ([]comment.Comment, error) {
	return nil, nil
}

func (s *Service) Delete(id int64) error {
	return nil
}
