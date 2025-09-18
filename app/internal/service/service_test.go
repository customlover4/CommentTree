package service

import (
	"CommentTree/internal/entities/comment"
	"CommentTree/internal/storage"
	"errors"
	"testing"
)

type StorageMock struct {
	createF func(c comment.Comment) (int64, error)
	getF    func(parentID int64, opts *comment.GetterOpts) (*comment.CommentView, error)
	deleteF func(id int64) error
}

func (sm *StorageMock) CreateComment(c comment.Comment) (int64, error) {
	return sm.createF(c)
}

func (sm *StorageMock) Comments(parentID int64, opts *comment.GetterOpts) (*comment.CommentView, error) {
	return sm.getF(parentID, opts)
}

func (sm *StorageMock) DeleteComment(id int64) error {
	return sm.deleteF(id)
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		str str
		// Named input parameters for target function.
		c    comment.Comment
		want error
	}{
		{
			name: "good",
			str: &StorageMock{
				createF: func(c comment.Comment) (int64, error) {
					return 0, nil
				},
			},
			c: comment.Comment{
				Message: "hi",
			},
			want: nil,
		},
		{
			name: "wrong comment to create",
			str: &StorageMock{
				createF: func(c comment.Comment) (int64, error) {
					return 0, nil
				},
			},
			c: comment.Comment{
				Message: "",
			},
			want: ErrWrongData,
		},
		{
			name: "wrong parent id",
			str: &StorageMock{
				createF: func(c comment.Comment) (int64, error) {
					return 0, nil
				},
			},
			c: comment.Comment{
				Message:  "asd",
				ParentID: -1,
			},
			want: ErrWrongData,
		},
		{
			name: "violates foreign key",
			str: &StorageMock{
				createF: func(c comment.Comment) (int64, error) {
					return 0, storage.ErrWrongForeignKey
				},
			},
			c: comment.Comment{
				Message:  "asd",
				ParentID: 123,
			},
			want: ErrWrongData,
		},
		{
			name: "internal error storage",
			str: &StorageMock{
				createF: func(c comment.Comment) (int64, error) {
					return 0, errors.New("unknown error")
				},
			},
			c: comment.Comment{
				Message: "asd",
			},
			want: ErrStorageInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.str)
			_, gotErr := s.CreateComment(tt.c)
			if !errors.Is(gotErr, tt.want) {
				t.Errorf("Create() want = %v, get %v", tt.want, gotErr)
			}
		})
	}
}

func TestService_Comments(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		str str
		// Named input parameters for target function.
		parentID int64
		opts     *comment.GetterOpts

		want error
	}{
		{
			name: "good",
			str: &StorageMock{
				getF: func(parentID int64, opts *comment.GetterOpts) (*comment.CommentView, error) {
					return nil, nil
				},
			},
			parentID: 1,
			opts:     &comment.GetterOpts{},
			want:     nil,
		},
		{
			name: "wrong parentID",
			str: &StorageMock{
				getF: func(parentID int64, opts *comment.GetterOpts) (*comment.CommentView, error) {
					return nil, nil
				},
			},
			parentID: -1,
			opts:     &comment.GetterOpts{},
			want:     ErrWrongData,
		},
		{
			name: "not found",
			str: &StorageMock{
				getF: func(parentID int64, opts *comment.GetterOpts) (*comment.CommentView, error) {
					return nil, storage.ErrNotFound
				},
			},
			parentID: 1,
			opts:     &comment.GetterOpts{},
			want:     ErrNotFound,
		},
		{
			name: "unknown error",
			str: &StorageMock{
				getF: func(parentID int64, opts *comment.GetterOpts) (*comment.CommentView, error) {
					return nil, errors.New("test")
				},
			},
			parentID: 1,
			opts:     &comment.GetterOpts{},
			want:     ErrStorageInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.str)
			_, gotErr := s.Comments(tt.parentID, tt.opts)
			if !errors.Is(gotErr, tt.want) {
				t.Errorf("Comments() want = %v, get %v", tt.want, gotErr)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		str str
		// Named input parameters for target function.
		ID int64

		want error
	}{
		{
			name: "good",
			str: &StorageMock{
				deleteF: func(id int64) error {
					return nil
				},
			},
			ID:   1,
			want: nil,
		},
		{
			name: "wrong id",
			str: &StorageMock{
				deleteF: func(id int64) error {
					return nil
				},
			},
			ID:   0,
			want: ErrWrongData,
		},
		{
			name: "not affected",
			str: &StorageMock{
				deleteF: func(id int64) error {
					return storage.ErrNotAffected
				},
			},
			ID:   1,
			want: ErrNotAffected,
		},
		{
			name: "unknown error",
			str: &StorageMock{
				deleteF: func(id int64) error {
					return errors.New("test")
				},
			},
			ID:   1,
			want: ErrStorageInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.str)
			gotErr := s.DeleteComment(tt.ID)
			if !errors.Is(gotErr, tt.want) {
				t.Errorf("Delete() want = %v, get %v", tt.want, gotErr)
			}
		})
	}
}
