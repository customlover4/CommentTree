package request

type CreateComment struct {
	Message  string `json:"message"`
	ParentID int64  `json:"parent_id"`
}

func (cc *CreateComment) Validate() string {
	if cc.Message == "" {
		return "empty message"
	}
	if cc.ParentID < 0 {
		return "wrong parent id"
	}

	return ""
}

type DeleteComment struct {
	ID int64 `json:"id"`
}

func (dc *DeleteComment) Validate() string {
	if dc.ID <= 0 {
		return "wrong id"
	}

	return ""
}
