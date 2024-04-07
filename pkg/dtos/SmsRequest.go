package dtos

type PostSenderIDRequest struct {
	SenderCode string `json:"sender_id"`
	Type       string `json:"type"`
	Language   string `json:"language"`
}

type GetSenderIDRequest struct {
	SenderCode string `form:"sender_id"`
	Type       string `form:"type"`
	UserID     string `form:"user_id"`
	Status     string `form:"status"`
}

type ApproveSenderIDRequest struct {
	SenderCode string `json:"sender_id"`
	UserID     string `json:"user_id"`
}

type DeleteSenderIDRequest struct {
	SenderCode string `json:"sender_id"`
}
