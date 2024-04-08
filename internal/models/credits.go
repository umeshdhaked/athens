package models

const (
	TableCredits = "Credits"
)

const (
	// Columns
	ColumnCreditsID          = "ID"
	ColumnCreditsUserID      = "UserID"
	ColumnCreditsCredits     = "Credits"
	ColumnCreditsCreditsLeft = "CreditsLeft"
)

type Credits struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Credits     float64 `json:"credits"`
	CreditsLeft float64 `json:"credits_left"`
	BaseModel
}
