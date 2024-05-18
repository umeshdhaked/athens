package models

const (
	TableSmsCampaign = "sms_campaign"
)

const (
	// States
	SmsCampaignStateCreated     = "CREATED"
	SmsCampaignStateDeActivated = "DEACTIVATED"
	SmsCampaignStateInProgress  = "IN_PROGRESS"
	SmsCampaignStateExecuted    = "EXECUTED"
)

type SmsCampaign struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ScheduledAt int    `json:"scheduled_at"`
	Status      string `json:"status"`
	UserID      string `json:"user_id"`
	TemplateID  int64  `json:"template_id"`
	SenderCode  string `json:"sender_code"`
	GroupName   string `json:"group_name"`
	Type        string `json:"type"`
	BaseModel
}

func (o *SmsCampaign) TableName() string {
	return TableSmsCampaign
}

func (o *SmsCampaign) GetID() int64 {
	return o.ID
}

func (o *SmsCampaign) SetID(id int64) {
	o.ID = id
}
