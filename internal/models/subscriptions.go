package models

const (
	TablePricing          = "Pricing"
	TableUserSubscription = "UserSubscriptions"
	TableCreditAudits     = "CreditsAudit"
)

const (
	ColumnId           = "Id"
	ColumnUserId       = "UserId"
	ColumnType         = "Type"
	ColumnPricingState = "PricingState"
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

type UserSubscription struct {
	Id        string `json:"id"`
	PricingId string `json:"pricing_id"`
	UserId    string `json:"user_id"`
	Type      string `json:"type"`
	SubType   string `json:"sub_type"`
	SubStatus string `json:"status"`
	AddedBy   string `json:"added_by"`
	BaseModel
}

type CreditAudits struct {
	Id             string  `json:"id"`
	Category       string  `json:"category"`
	SubCategory    string  `json:"sub_category"`
	DeductedAmount float64 `json:"deducted_amount"`
	AddedAmount    float64 `json:"added_amount"`
	CreditId       string  `json:"credit_id"`
	UserId         string  `json:"user_id"`
	PaymentOrderId string  `json:"payment_order_id"`
	BaseModel
}
