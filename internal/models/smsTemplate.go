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
	ColumnSmsTemplateID           = "id"
	ColumnSmsTemplateUserID       = "user_id"
	ColumnSmsTemplateSenderID     = "sender_id"
	ColumnSmsTemplateSenderCode   = "sender_coe"
	ColumnSmsTemplateTemplateCode = "template_code"
	ColumnSmsTemplateBody         = "body"
	ColumnSmsTemplateStatus       = "status"
	ColumnSmsTemplateType         = "type"
	ColumnSmsTemplateLanguage     = "language"
)

type SmsTemplate struct {
	ID           string `json:"id" dynamodbav:"id"`
	UserID       string `json:"user_id" dynamodbav:"user_id"`
	SenderID     string `json:"sender_id" dynamodbav:"sender_id"`
	SenderCode   string `json:"sender_code" dynamodbav:"sender_code"`
	TemplateCode string `json:"template_code" dynamodbav:"template_code"`
	Body         string `json:"body" dynamodbav:"body"`
	Status       string `json:"status" dynamodbav:"status"`
	Type         string `json:"type" dynamodbav:"type"`
	Length       int    `json:"length" dynamodbav:"length"`
	Language     string `json:"language" dynamodbav:"language"`
	BaseModel
}
