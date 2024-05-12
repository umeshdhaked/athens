package models

const (
	TableSmsAudit = "SmsAudit"
)

const (
	// Statues
	SmsAuditStatusDelivered = "DELIVERED"

	// Triggered Modes
	ModeSmsAuditInstant = "Instant"
	ModeSmsAuditBulk    = "Bulk"

	// Columns
	ColumnSmsAuditID              = "ID"
	ColumnSmsAuditUserID          = "UserID"
	ColumnSmsAuditCreditsConsumed = "CreditsConsumed"
	ColumnSmsAuditTemplateID      = "TemplateCode"
	ColumnSmsAuditSenderCode      = "SenderCode"
	ColumnSmsAuditContactID       = "ContactID"
	ColumnSmsAuditStatus          = "Status"
	ColumnSmsAuditTriggeredMode   = "TriggeredMode"
)

type SmsAudit struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	CreditsConsumed float64 `json:"credits_consumed"`
	TemplateID      string  `json:"template_id"`
	SenderCode      string  `json:"sender_code"`
	ContactID       string  `json:"contact_id"`
	Status          string  `json:"status"`
	TriggeredMode   string  `json:"triggered_mode"`
	BaseModel
}
