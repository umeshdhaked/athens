package models

const (
	TablePricing = "pricing"
)

type Pricing struct {
	ID          int64   `json:"id"`
	Category    string  `json:"category"`
	SubCategory string  `json:"sub_category"`
	PricingType string  `json:"pricing_type"`
	Rates       float64 `json:"rates"`
	State       string  `json:"state"`
	BaseModel
}

func (o *Pricing) TableName() string {
	return TablePricing
}

func (o *Pricing) GetID() int64 {
	return o.ID
}

func (o *Pricing) SetID(id int64) {
	o.ID = id
}
