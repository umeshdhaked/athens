package models

const (
	TableSmsSender = "SmsSender"

	// States
	SmsSenderStateCreated     = "CREATED"
	SmsSenderStateApproved    = "APPROVED"
	SmsSenderStateDeActivated = "DEACTIVATED"
)

const (
	ColumnSmsSenderID     = "ID"
	ColumnSmsSenderCode   = "Code"
	ColumnSmsSenderUserID = "UserID"
	ColumnSmsSenderType   = "Type"
	ColumnSmsSenderStatus = "Status"
)

type SmsSender struct {
	ID     string `json:"id"`
	Code   string `json:"name"`
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	BaseModel
}
