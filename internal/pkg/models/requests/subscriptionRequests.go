package requests

type PricingRequest struct {
	Category   string
	SubCatgory string
	Type       string
	Rates      float32
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
