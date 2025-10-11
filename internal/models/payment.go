package models

import "time"

const (
	CurrencyINR       = "INR"
	PaymentOrgName    = "umeshdhaked Solutions"
	TestPaymentKey    = "rzp_test_NmbdVYr5EuOOo5"
	TestPaymentSecret = "vW4mQzdP5SL0VlSrRNlkLc7Y"
	TableInvoices     = "Invoice"

	TablePayments           = "payments"
	SQLColumnInvoiceOrderId = "order_id"
)

type Payments struct {
	Amount         float64 `json:"amount"`
	AmountPaid     float64 `json:"amount_paid"`
	Mobile         string  `json:"mobile"`
	UserId         float64 `json:"user_id"`
	OrderCreatedAt float64 `json:"order_created_at"`
	AmountDue      float64 `json:"amount_due"`
	Currency       string  `json:"currency"`
	Receipt        string  `json:"receipt"`
	OrderId        string  `json:"order_id"`
	OfferId        string  `json:"offer_id"`
	Entity         string  `json:"entity"`
	Attempts       float64 `json:"attempts"`
	Status         string  `json:"status"`
	Id             int64   `json:"id"`
	BaseModel
}

func (o *Payments) TableName() string {
	return TablePayments
}

func (o *Payments) GetID() int64 {
	return o.Id
}

func (o *Payments) SetID(id int64) {
	o.Id = id
}

func (o *Payments) PopulateFromMap(body map[string]interface{}) {
	if nil != body["entity"] {
		o.Entity = body["entity"].(string)
	}
	if nil != body["amount"] {
		o.Amount = body["amount"].(float64)
	}
	if nil != body["amount_paid"] {
		o.AmountPaid = body["amount_paid"].(float64)
	}
	if nil != body["receipt"] {
		o.Receipt = body["receipt"].(string)
	}
	if nil != body["offer_id"] {
		o.OfferId = body["offer_id"].(string)
	}
	if nil != body["status"] {
		o.Status = body["status"].(string)
	}
	if nil != body["notes"] && nil != (body["notes"].(map[string]interface{}))["mobile"] {
		o.Mobile = (body["notes"].(map[string]interface{}))["mobile"].(string)
	}
	if nil != body["notes"] && nil != (body["notes"].(map[string]interface{}))["userId"] {
		o.UserId = (body["notes"].(map[string]interface{}))["userId"].(float64)
	}
	if nil != body["created_at"] {
		o.OrderCreatedAt = body["created_at"].(float64)
	}
	if nil != body["id"] {
		o.OrderId = body["id"].(string)
	}
	if nil != body["amount_due"] {
		o.AmountDue = body["amount_due"].(float64)
	}
	if nil != body["currency"] {
		o.Currency = body["currency"].(string)
	}
	if nil != body["attempts"] {
		o.Attempts = body["attempts"].(float64)
	}
	o.UpdatedAt = time.Now().Unix()
	if o.BaseModel.CreatedAt == 0 {
		o.CreatedAt = time.Now().Unix()
	}
}

type Invoice struct {
	ID      int64
	OrderId string
	Status  string
	UserId  float64
	Receipt string
	BaseModel
}

func (o *Invoice) TableName() string {
	return TableInvoices
}

func (o *Invoice) GetID() int64 {
	return o.ID
}

func (o *Invoice) SetID(id int64) {
	o.ID = id
}
