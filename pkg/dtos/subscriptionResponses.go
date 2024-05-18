package dtos

type PricingResponse struct {
	Id         int64
	Category   string
	SubCatgory string
	Type       string
	Rates      float64
	Status     string
}

type SubscriptionResponse struct {
	Id        int64
	PricingId int64
	UserId    int64
	Type      string
	SubType   string
	Status    string
	AddedBy   string
	CreatedAt int64
	DeletedAt int64
}

type CreditsResponse struct {
	Id              int64
	UserMobile      string
	InitialCredit   float64
	RemainingCredit float64
	CreatedAt       int64
}
