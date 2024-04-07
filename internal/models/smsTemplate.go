package models

const (
	TableSmsTemplate = "SmsTemplate"
)

const (
	// States
	SmsTemplateStateCreated     = "CREATED"
	SmsTemplateStateApproved    = "APPROVED"
	SmsTemplateStateDeActivated = "DEACTIVATED"

	// Columns
	ColumnSmsTemplateID         = "ID"
	ColumnSmsTemplateUserID     = "UserID"
	ColumnSmsTemplateSenderID   = "SenderID"
	ColumnSmsTemplateSenderCode = "SenderCode"
	ColumnSmsTemplateTemplateID = "TemplateID"
	ColumnSmsTemplateBody       = "Body"
	ColumnSmsTemplateStatus     = "Status"
	ColumnSmsTemplateType       = "Type"
	ColumnSmsTemplateLanguage   = "Language"
)

type SmsTemplate struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	SenderID   string `json:"sender_id"`
	SenderCode string `json:"sender_code"`
	TemplateID string `json:"template_id"`
	Body       string `json:"body"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	Language   string `json:"language"`
	BaseModel
}
