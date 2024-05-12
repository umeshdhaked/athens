package models

const (
	TablePricing = "Pricing"
)

const (
	// Columns
	ColumnPricingID = "Id"
)

type Pricing struct {
	Id           string  `json:"id"`
	Category     string  `json:"category"`
	SubCategory  string  `json:"sub_category"`
	PricingType  string  `json:"pricing_type"`
	Rates        float64 `json:"rates"`
	PricingState string  `json:"pricing_state"`
	BaseModel
}
