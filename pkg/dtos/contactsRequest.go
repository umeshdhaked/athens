package dtos

type GetContactsRequest struct {
	Name      string `form:"name"`
	Email     string `form:"email"`
	Mobile    string `form:"mobile"`
	GroupName string `form:"group_name"`
	From      int    `form:"from"`
	To        int    `form:"to"`
	Limit     int    `form:"limit"`
}

type GetGroupContactsRequest struct {
	UserID    string `json:"user_id" form:"user_id"`
	GroupName string `json:"group_name" form:"group_name"`
}
