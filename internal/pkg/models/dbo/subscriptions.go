package dbo

type Pricing struct {
	Id           string  `json:"id"`
	Category     string  `json:"category"`
	SubCatgory   string  `json:"sub_category"`
	PricingType  string  `json:"pricing_type"`
	Rates        float32 `json:"rates"`
	PricingState string  `json:"pricing_state"`
	CreatedAt    int64   `json:"createdAt"`
	DeletedAt    int64   `json:"deletedAt"`
}

type UserSubscription struct {
	Id        string `json:"id"`
	PricingId string `json:"pricing_id"`
	UserId    string `json:"user_id"`
	Type      string `json:"type"`
	SubType   string `json:"sub_type"`
	Status    string `json:"status"`
	AddedBy   string `json:"added_by"`
	CreatedAt int64  `json:"createdAt"`
	DeletedAt int64  `json:"deletedAt"`
}

type Credits struct {
	Id              string  `json:"id"`
	UserId          string  `json:"user_id"`
	InitialCredit   float32 `json:"initial_credit"`
	RemainingCredit float32 `json:"remaining_credit"`
	CreatedAt       int64   `json:"createdAt"`
	DeletedAt       int64   `json:"deletedAt"`
}

type CreditAudits struct {
	Id            string  `json:"id"`
	Category      string  `json:"category"`
	SubCategory   string  `json:"sub_category"`
	DeductedAmout float32 `json:"deducted_amount"`
	AddedAmount   float32 `json:"added_amount"`
	CreditId      string  `json:"credit_id"`
	UserId        string  `json:"user_id"`
	UpdatedAt     int64   `json:"updated_at"`
}
