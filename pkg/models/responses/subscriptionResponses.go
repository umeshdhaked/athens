package responses

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
