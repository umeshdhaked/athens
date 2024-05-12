package models

const (
	TableSmsCampaign = "SmsCampaign"
)

const (
	// States
	SmsCampaignStateCreated     = "CREATED"
	SmsCampaignStateDeActivated = "DEACTIVATED"
	SmsCampaignStateInProgress  = "IN_PROGRESS"
	SmsCampaignStateExecuted    = "EXECUTED"

	// Columns
	ColumnSmsCampaignID          = "ID"
	ColumnSmsCampaignName        = "Name"
	ColumnSmsCampaignScheduledAt = "ScheduledAt"
	ColumnSmsCampaignStatus      = "Status"
	ColumnSmsCampaignUserID      = "UserID"
	ColumnSmsCampaignTemplateID  = "TemplateCode"
	ColumnSmsCampaignSenderID    = "SenderID"
	ColumnSmsCampaignType        = "Type"
)

type SmsCampaign struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ScheduledAt int    `json:"scheduled_at"`
	Status      string `json:"status"`
	UserID      string `json:"user_id"`
	TemplateID  string `json:"template_id"`
	SenderCode  string `json:"sender_code"`
	GroupName   string `json:"group_name"`
	Type        string `json:"type"`
	BaseModel
}
