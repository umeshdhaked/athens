package models

const (
	TableSmsAudit = "sms_audit"
)

const (
	// Statues
	SmsAuditStatusDelivered = "DELIVERED"

	// Triggered Modes
	ModeSmsAuditInstant = "Instant"
	ModeSmsAuditBulk    = "Bulk"
)

type SmsAudit struct {
	ID              int64   `json:"id"`
	UserID          string  `json:"user_id"`
	CreditsConsumed float64 `json:"credits_consumed"`
	TemplateID      int64   `json:"template_id"`
	SenderCode      string  `json:"sender_code"`
	ContactID       int64   `json:"contact_id"`
	Status          string  `json:"status"`
	TriggeredMode   string  `json:"triggered_mode"`
	BaseModel
}

func (o *SmsAudit) TableName() string {
	return TableSmsAudit
}

func (o *SmsAudit) GetID() int64 {
	return o.ID
}

func (o *SmsAudit) SetID(id int64) {
	o.ID = id
}
