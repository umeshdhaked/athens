package models

const (
	// Columns
	ColumnLocksKey                 = "key"
	ColumnLocksData                = "data"
	ColumnLocksOwnerName           = "ownerName"
	ColumnLocksRecordVersionNumber = "recordVersionNumber"
	ColumnLocksIsReleased          = "isReleased"
	ColumnLocksLeaseDuration       = "leaseDuration"
)

type Locks struct {
	Key                 string  `json:"key"`
	Data                string  `json:"data"`
	OwnerName           float64 `json:"ownerName"`
	RecordVersionNumber string  `json:"recordVersionNumber"`
	IsReleased          string  `json:"isReleased"`
	LeaseDuration       string  `json:"leaseDuration"`
}
