package request

import (
	"testing"
)

func TestCreateComment_Validate(t *testing.T) {
	tests := []struct {
		name     string // description of this test case
		message  string
		parentID int64
		want  bool
	}{
		{
			name:     "good",
			message:  "hi",
			parentID: 0,
			want:     false,
		},
		{
			name:     "bad message",
			message:  "",
			parentID: 0,
			want:     true,
		},
		{
			name:     "bad parent id",
			message:  "123",
			parentID: -1,
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cc CreateComment
			cc.Message = tt.message
			cc.ParentID = tt.parentID
			got := cc.Validate()
			if tt.want && got == "" {
				t.Errorf("Validate() = %v, want %T", got, tt.want)
			} else if !tt.want && got != "" {
				t.Errorf("Validate() = %v, want %T", got, tt.want)
			}
		})
	}
}
