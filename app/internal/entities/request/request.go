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

type GetComments struct {
	ParentID     int64  `form:"parent"`
	Substr       string `form:"substr"`
	Page         int    `form:"page"`
	SearchGlobal bool   `form:"search_global"`
}

func (gc *GetComments) Validate() string {
	if gc.ParentID < 0 {
		return "wrong id, id should be >= 0"
	}
	if gc.Page < 0 {
		return "wrong page, page should be >= 0"
	}

	if gc.Page == 0 {
		gc.Page++
	}

	return ""
}
