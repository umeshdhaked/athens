package models

const (
	TableSmsSender = "sms_sender"

	// States
	SmsSenderStateCreated     = "CREATED"
	SmsSenderStateApproved    = "APPROVED"
	SmsSenderStateDeActivated = "DEACTIVATED"
)

type SmsSender struct {
	ID     int64  `json:"id"`
	Code   string `json:"name"`
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	BaseModel
}

func (s *SmsSender) TableName() string {
	return TableSmsSender
}

func (s *SmsSender) GetID() int64 {
	return s.ID
}

func (s *SmsSender) SetID(id int64) {
	s.ID = id
}
