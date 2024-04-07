package dtos

type PostSenderCodeRequest struct {
	SenderCode string `json:"sender_code"`
	Type       string `json:"type"`
	Language   string `json:"language"`
}

type GetSenderCodeRequest struct {
	SenderCode string `form:"sender_code"`
	Type       string `form:"type"`
	UserID     string `form:"user_id"`
	Status     string `form:"status"`
}

type ApproveSenderCodeRequest struct {
	SenderCode string `json:"sender_code"`
	UserID     string `json:"user_id"`
}

type DeleteSenderCodeRequest struct {
	SenderCode string `json:"sender_code"`
}

type PostSmsTemplateRequest struct {
	SenderID   string `json:"sender_id"`
	SenderCode string `json:"sender_code"`
	TemplateID string `json:"template_id"`
	Body       string `json:"body"`
	Type       string `json:"type"`
	Language   string `json:"language"`
}

type GetSmsTemplateRequest struct {
	UserID     string `form:"user_id"`
	SenderCode string `form:"sender_code"`
	TemplateID string `form:"template_id"`
	Status     string `form:"status"`
}

type ApproveSmsTemplateRequest struct {
	TemplateID string `json:"template_id"`
	UserID     string `json:"user_id"`
}

type UpdateSmsTemplateRequest struct {
	TemplateID string `json:"template_id"`
	UserID     string `json:"user_id"`
	Body       string `json:"body"`
	Status     string `json:"status"`
}

type DeActivateSmsTemplateRequest struct {
	TemplateID string `json:"template_id"`
	UserID     string `json:"user_id"`
}
