package dtos

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

type DeactivateSubscriptionRequest struct {
	Id string
}

type AddCreditsRequest struct {
	UserMobile    string
	InitialCredit float32
}

type ChargeUserRequest struct {
	UserId      string
	Category    string
	SubCategory string
	UnitCount   float32
}
