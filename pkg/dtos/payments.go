package dtos

type PaymentOrderRequest struct {
	Amount int64 `json:"amount"`
}

type PaymentOrderResponse struct {
	Key        string  `json:"key"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	OrgName    string  `json:"orgName"`
	RzpOrderId string  `json:"rzpOrderId"`
}

type UpdatePaymentOrderRequest struct {
	RazorpayPaymentId string `json:"razorpayPaymentId"`
	RazorpayOrderId   string `json:"razorpayOrderId"`
	RazorpaySignature string `json:"razorpaySignature"`
	ErrorCode         string `json:"errorCode"`
	ErrorDescription  string `json:"errorDescription"`
	ErrorSource       string `json:"errorSource"`
	ErrorStep         string `json:"errorStep"`
	ErrorReason       string `json:"errorReason"`
}

type UpdatePaymentResponse struct {
	OrderStatus   string `json:"orderStatus"`
	PaymentStatus string `json:"paymentStatus"`
}

type PaymentWebhookRequest struct {
	Entity    string                                  `json:"entity"`
	AccountId string                                  `json:"account_id"`
	Event     string                                  `json:"event"`
	Contains  []string                                `json:"contains"`
	CreatedAt int64                                   `json:"created_at"`
	Payload   map[string]PaymentWebhookRequestPayload `json:"payload"`
}

type PaymentWebhookRequestPayload struct {
	Entity map[string]interface{} `json:"entity"`
}

type PaymentWebhookResponse struct {
	Status string `json:"status"`
}

type GetPaymentStatusRequest struct {
	OrderId string `json:"orderId"`
}

type GetPaymentStatusResponse struct {
	Status string `json:"status"`
}
