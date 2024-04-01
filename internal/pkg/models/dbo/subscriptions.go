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
	Id              string `json:"id"`
	UserId          string `json:"userId"`
	InitialCredit   string `json:"initialCredit"`
	RemainingCredit string `json:"remainingCredit"`
	CreatedAt       string `json:"createdAt"`
	DeletedAt       string `json:"deletedAt"`
}
