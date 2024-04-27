package dtos

type GetS3ProcessingRequest struct {
	Name   string `form:"name"`
	Status string `form:"status"`
	Type   string `form:"type"`
	From   int    `form:"from"`
	To     int    `form:"to"`
}
