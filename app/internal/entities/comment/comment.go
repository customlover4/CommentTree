package comment

const (
	// Elements per page for pagination.
	PageElements = 10
)

type Comment struct {
	ID       int64
	Message  string
	ParentID int64
	HaveNext bool `json:"have_next"`
}

type CommentView struct {
	Parent Comment   `json:"parent"`
	Childs []Comment `json:"childs"`
}

type GetterOpts struct {
	// For searching in all comments by substr.
	SearchGlobal bool

	Page   int
	Substr string
}

func (g *GetterOpts) Empty() bool {
	if g.Page == 0 && g.Substr == "" {
		return true
	}

	return false
}
