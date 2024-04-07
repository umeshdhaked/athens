package dtos

type PricingResponse struct {
	Id         string
	Category   string
	SubCatgory string
	Type       string
	Rates      float32
	Status     string
}

type SubscriptionResponse struct {
	Id        string
	PricingId string
	UserId    string
	Type      string
	SubType   string
	Status    string
	AddedBy   string
	CreatedAt int64
	DeletedAt int64
}

type CreditsResponse struct {
	Id              string
	UserMobile      string
	InitialCredit   float32
	RemainingCredit float32
	CreatedAt       int64
}
