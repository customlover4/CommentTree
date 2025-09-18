package comment

type Comment struct {
	ID       int64
	Message  string
	ParentID int64
}

type CommentView struct {
	Parent Comment   `json:"parent"`
	Childs []Comment `json:"childs"`
}

type GetterOpts struct {
	Page   int    `form:"page"`
	Substr string `form:"substr"`
}

func (g *GetterOpts) Empty() bool {
	if g.Page == 0 && g.Substr == "" {
		return true
	}

	return false
}
