package postgres

import "CommentTree/internal/entities/comment"

func (p *Postgres) Create(c comment.Comment) (int64, error) {
	return 0, nil
}

func (p *Postgres) Comments(id int64) ([]comment.Comment, error) {
	return nil, nil
}

func (p *Postgres) Delete(id int64) error {
	return nil
}
