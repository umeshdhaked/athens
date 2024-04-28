package dtos

type PricingRequest struct {
	Category   string
	SubCatgory string
	Type       string
	Rates      float64
}

type DeactivatePricingRequest struct {
	Id     string
	Enable bool
}

type UserDefaultSubscriptionRequest struct {
	UserMobile string
}

type UserSubscriptionRequest struct {
	UserDefaultSubscriptionRequest
	PricingId string
}

type FetchSubscriptionRequest struct {
	UserMobile string
}

type DeactivateSubscriptionRequest struct {
	Id     string
	Enable bool
}

type AddCreditsRequest struct {
	UserMobile     string
	InitialCredit  float64
	PaymentOrderId string
}

type ChargeUserRequest struct {
	UserId      string
	Category    string
	SubCategory string
	UnitCount   float64
}
