package storage

type Storage struct {
	Comment *CommentStorage
}

func New(cs *CommentStorage) *Storage {
	return &Storage{
		Comment: cs,
	}
}