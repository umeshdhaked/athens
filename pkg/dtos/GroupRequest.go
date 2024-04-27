package dtos

type UploadGroupContactsRequest struct {
	Name string `json:"name" form:"name"`
}

type GetGroupRequest struct {
	UserID string `json:"user_id" form:"user_id"`
	Name   string `json:"name" form:"name"`
	From   int    `form:"from"`
	To     int    `form:"to"`
}
