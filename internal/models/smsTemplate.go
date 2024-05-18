package models

const (
	TableSmsTemplate = "sms_template"
)

const (
	// States
	SmsTemplateStateCreated     = "CREATED"
	SmsTemplateStateApproved    = "APPROVED"
	SmsTemplateStateDeActivated = "DEACTIVATED"
)

type SmsTemplate struct {
	ID           int64  `json:"id" dynamodbav:"id"`
	UserID       string `json:"user_id" dynamodbav:"user_id"`
	SenderID     int64  `json:"sender_id" dynamodbav:"sender_id"`
	SenderCode   string `json:"sender_code" dynamodbav:"sender_code"`
	TemplateCode string `json:"template_code" dynamodbav:"template_code"`
	Body         string `json:"body" dynamodbav:"body"`
	Status       string `json:"status" dynamodbav:"status"`
	Type         string `json:"type" dynamodbav:"type"`
	Length       int    `json:"length" dynamodbav:"length"`
	Language     string `json:"language" dynamodbav:"language"`
	BaseModel
}

func (o *SmsTemplate) TableName() string {
	return TableSmsTemplate
}

func (o *SmsTemplate) GetID() int64 {
	return o.ID
}

func (o *SmsTemplate) SetID(id int64) {
	o.ID = id
}
