package dtos

type GetGroupContactsRequest struct {
	Name string `json:"name" form:"name"`
}

type UploadGroupContactsRequest struct {
	Name string `json:"name" form:"name"`
}
