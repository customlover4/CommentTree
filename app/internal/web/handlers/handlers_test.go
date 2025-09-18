package handlers

import (
	"CommentTree/internal/entities/comment"
	"CommentTree/internal/entities/request"
	"CommentTree/internal/service"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type ServiceMock struct {
	createF func(c comment.Comment) (int64, error)
	getF    func(parentID int64, opts *comment.GetterOpts) (*comment.CommentView, error)
	deleteF func(id int64) error
}

func (sm *ServiceMock) CreateComment(c comment.Comment) (int64, error) {
	return sm.createF(c)
}

func (sm *ServiceMock) Comments(parentID int64, opts *comment.GetterOpts) (*comment.CommentView, error) {
	return sm.getF(parentID, opts)
}

func (sm *ServiceMock) DeleteComment(id int64) error {
	return sm.deleteF(id)
}

func TestCreateComment(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s       servicer
		body    request.CreateComment // set struct and test automaticly marshal it in json
		bodyStr string                // if we want set json in string format
		want    int
	}{
		{
			name: "good",
			s: &ServiceMock{
				createF: func(c comment.Comment) (int64, error) {
					return 0, nil
				},
			},
			body: request.CreateComment{
				Message:  "hihi",
				ParentID: 12,
			},
			want: http.StatusOK,
		},
		{
			name: "bad json",
			s: &ServiceMock{
				createF: func(c comment.Comment) (int64, error) {
					return 0, nil
				},
			},
			bodyStr: `{"}`,
			want:    http.StatusBadRequest,
		},
		{
			name: "bad data in request",
			s: &ServiceMock{
				createF: func(c comment.Comment) (int64, error) {
					return 0, nil
				},
			},
			body: request.CreateComment{
				Message:  "",
				ParentID: 12,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "service err: wrong data",
			s: &ServiceMock{
				createF: func(c comment.Comment) (int64, error) {
					return 0, service.ErrWrongData
				},
			},
			body: request.CreateComment{
				Message:  "hihi",
				ParentID: 12,
			},
			want: http.StatusServiceUnavailable,
		},
		{
			name: "service err: unknown err",
			s: &ServiceMock{
				createF: func(c comment.Comment) (int64, error) {
					return 0, errors.New("unknown")
				},
			},
			body: request.CreateComment{
				Message:  "hihi",
				ParentID: 12,
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/endpoint"

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST(url, CreateComment(tt.s))

			var bytesBody []byte
			if tt.bodyStr != "" {
				bytesBody = []byte(tt.bodyStr)
			} else {
				tmp, err := json.Marshal(&tt.body)
				if err != nil {
					t.Error(err.Error())
				}
				bytesBody = tmp
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodPost, url, bytes.NewReader(bytesBody),
			)

			router.ServeHTTP(rr, req)

			if rr.Result().StatusCode != tt.want {
				t.Errorf(
					"handler CreateComment() = %v, want %v",
					rr.Result().StatusCode, tt.want,
				)
			}
		})
	}
}

func TestDeleteComment(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s     servicer
		param string
		want  int
	}{
		{
			name: "good",
			s: &ServiceMock{
				deleteF: func(id int64) error {
					return nil
				},
			},
			param: "12",
			want:  http.StatusOK,
		},
		{
			name: "bad json",
			s: &ServiceMock{
				deleteF: func(id int64) error {
					return nil
				},
			},
			param: "2sdg",
			want:  http.StatusBadRequest,
		},
		{
			name: "not valid data in json",
			s: &ServiceMock{
				deleteF: func(id int64) error {
					return nil
				},
			},
			param: "-1",
			want:  http.StatusBadRequest,
		},
		{
			name: "not valid id",
			s: &ServiceMock{
				deleteF: func(id int64) error {
					return service.ErrWrongData
				},
			},
			param: "1",
			want:  http.StatusServiceUnavailable,
		},
		{
			name: "not affected",
			s: &ServiceMock{
				deleteF: func(id int64) error {
					return service.ErrNotAffected
				},
			},
			param: "1",
			want:  http.StatusServiceUnavailable,
		},
		{
			name: "unknown err",
			s: &ServiceMock{
				deleteF: func(id int64) error {
					return errors.New("unknown")
				},
			},
			param: "1",
			want:  http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/endpoint/"

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.DELETE(url+":id", DeleteComment(tt.s))

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(
				http.MethodDelete, url+tt.param, nil,
			)

			router.ServeHTTP(rr, req)

			if rr.Result().StatusCode != tt.want {
				t.Errorf(
					"handler CreateComment() = %v, want %v",
					rr.Result().StatusCode, tt.want,
				)
			}
		})
	}
}
