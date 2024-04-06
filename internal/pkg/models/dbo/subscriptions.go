package dbo

type Pricing struct {
	Id           string  `json:"id" dynamodbav:"id"`
	Category     string  `json:"category" dynamodbav:"category"`
	SubCatgory   string  `json:"sub_category" dynamodbav:"sub_category"`
	PricingType  string  `json:"pricing_type" dynamodbav:"pricing_type"`
	Rates        float32 `json:"rates" dynamodbav:"rates"`
	PricingState string  `json:"pricing_state" dynamodbav:"pricing_state"`
	CreatedAt    int64   `json:"createdAt" dynamodbav:"createdAt"`
	DeletedAt    int64   `json:"deletedAt" dynamodbav:"deletedAt"`
}

type UserSubscription struct {
	Id        string `json:"id" dynamodbav:"id"`
	PricingId string `json:"pricing_id" dynamodbav:"pricing_id"`
	UserId    string `json:"user_id" dynamodbav:"user_id"`
	Type      string `json:"type" dynamodbav:"type"`
	SubType   string `json:"sub_type" dynamodbav:"sub_type"`
	Status    string `json:"status" dynamodbav:"status"`
	AddedBy   string `json:"added_by" dynamodbav:"added_by"`
	CreatedAt int64  `json:"createdAt" dynamodbav:"createdAt"`
	DeletedAt int64  `json:"deletedAt" dynamodbav:"deletedAt"`
}

type Credits struct {
	Id              string  `json:"id" dynamodbav:"id"`
	UserId          string  `json:"user_id" dynamodbav:"user_id"`
	InitialCredit   float32 `json:"initial_credit" dynamodbav:"initial_credit"`
	RemainingCredit float32 `json:"remaining_credit" dynamodbav:"remaining_credit"`
	CreatedAt       int64   `json:"createdAt" dynamodbav:"createdAt"`
	DeletedAt       int64   `json:"deletedAt" dynamodbav:"deletedAt"`
}

type CreditAudits struct {
	Id            string  `json:"id" dynamodbav:"id"`
	Category      string  `json:"category" dynamodbav:"category"`
	SubCategory   string  `json:"sub_category" dynamodbav:"sub_category"`
	DeductedAmout float32 `json:"deducted_amount" dynamodbav:"deducted_amount"`
	AddedAmount   float32 `json:"added_amount" dynamodbav:"added_amount"`
	CreditId      string  `json:"credit_id" dynamodbav:"credit_id"`
	UserId        string  `json:"user_id" dynamodbav:"user_id"`
	UpdatedAt     int64   `json:"updated_at" dynamodbav:"updated_at"`
}
