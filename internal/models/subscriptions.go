package models

const (
	TableSubscription = "Subscription"
)

const (
	ColumnSubscriptionsID            = "Id"
	ColumnSubscriptionsUserId        = "UserId"
	ColumnSubscriptionsPricingStatus = "status"
	ColumnSubscriptionsType          = "type"
	ColumnSubscriptionsSubType       = "sub_type"
)

type Subscription struct {
	Id        string `json:"ID" dynamodbav:"ID"`
	PricingId string `json:"pricing_id" dynamodbav:"pricing_id"`
	UserId    string `json:"UserId" dynamodbav:"UserId"`
	Type      string `json:"type" dynamodbav:"type"`
	SubType   string `json:"sub_type" dynamodbav:"sub_type"`
	Status    string `json:"status" dynamodbav:"status"`
	AddedBy   string `json:"added_by" dynamodbav:"added_by"`
	BaseModel
}
