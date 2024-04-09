package models

type Pricing struct {
	Id           string  `json:"id" dynamodbav:"id"`
	Category     string  `json:"category" dynamodbav:"category"`
	SubCatgory   string  `json:"sub_category" dynamodbav:"sub_category"`
	PricingType  string  `json:"pricing_type" dynamodbav:"pricing_type"`
	Rates        float64 `json:"rates" dynamodbav:"rates"`
	PricingState string  `json:"pricing_state" dynamodbav:"pricing_state"`
	BaseModel
}

type UserSubscription struct {
	Id        string `json:"id" dynamodbav:"id"`
	PricingId string `json:"pricing_id" dynamodbav:"pricing_id"`
	UserId    string `json:"user_id" dynamodbav:"user_id"`
	Type      string `json:"type" dynamodbav:"type"`
	SubType   string `json:"sub_type" dynamodbav:"sub_type"`
	Status    string `json:"status" dynamodbav:"status"`
	AddedBy   string `json:"added_by" dynamodbav:"added_by"`
	BaseModel
}

type CreditAudits struct {
	Id            string  `json:"id" dynamodbav:"id"`
	Category      string  `json:"category" dynamodbav:"category"`
	SubCategory   string  `json:"sub_category" dynamodbav:"sub_category"`
	DeductedAmout float64 `json:"deducted_amount" dynamodbav:"deducted_amount"`
	AddedAmount   float64 `json:"added_amount" dynamodbav:"added_amount"`
	CreditId      string  `json:"credit_id" dynamodbav:"credit_id"`
	UserId        string  `json:"user_id" dynamodbav:"user_id"`
	BaseModel
}
