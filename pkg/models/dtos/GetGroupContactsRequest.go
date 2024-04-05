package dtos

type GetGroupContactsRequest struct {
	Name string `json:"name" query:"name"`
}
