package models

const (
	TableCreditAudits = "CreditsAudit"
)

type CreditAudits struct {
	Id             string  `json:"id"`
	Category       string  `json:"category"`
	SubCategory    string  `json:"sub_category"`
	DeductedAmount float64 `json:"deducted_amount"`
	AddedAmount    float64 `json:"added_amount"`
	CreditId       string  `json:"credit_id"`
	UserId         string  `json:"user_id"`
	PaymentOrderId string  `json:"payment_order_id"`
	BaseModel
}
