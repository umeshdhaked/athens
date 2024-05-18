package models

const (
	TableSubscription = "subscription"
)

type Subscription struct {
	ID        int64  `json:"ID" dynamodbav:"ID"`
	PricingId int64  `json:"pricing_id" dynamodbav:"pricing_id"`
	UserId    int64  `json:"UserId" dynamodbav:"UserId"`
	Type      string `json:"type" dynamodbav:"type"`
	SubType   string `json:"sub_type" dynamodbav:"sub_type"`
	Status    string `json:"status" dynamodbav:"status"`
	AddedBy   string `json:"added_by" dynamodbav:"added_by"`
	BaseModel
}

func (o *Subscription) TableName() string {
	return TableSubscription
}

func (o *Subscription) GetID() int64 {
	return o.ID
}

func (o *Subscription) SetID(id int64) {
	o.ID = id
}
