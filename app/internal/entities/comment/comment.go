package comment

type Comment struct {
	ID       int64
	Message  string
	ParentID int64
}

type GetterOpts struct {
	Page   int    `form:"page"`
	Substr string `form:"substr"`
}
