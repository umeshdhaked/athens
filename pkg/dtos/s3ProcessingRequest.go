package dtos

type GetCronProcessingRequest struct {
	Name   string `form:"name"`
	Status string `form:"status"`
	From   int    `form:"from"`
	To     int    `form:"to"`
}
