package models

const (
	CurrencyINR       = "INR"
	PaymentOrgName    = "FastBizTech Solutions"
	TestPaymentKey    = "rzp_test_NmbdVYr5EuOOo5"
	TestPaymentSecret = "vW4mQzdP5SL0VlSrRNlkLc7Y"
	TableRzpOrders    = "RzpOrders"
	TableRzpPayments  = "RzpPayments"
	TableRzpRefunds   = "RzpRefunds"
	TableFBTPayment   = "FBTPayment"

	ColumnOrderId     = "id"
	ColumnOrderStatus = "status"
)

type RzpOrder struct {
	Amount     int `dynamodbav:"amount"`
	AmountPaid int `dynamodbav:"amount_paid"`
	Notes      struct {
		Mobile string `dynamodbav:"mobile"`
		UserID string `dynamodbav:"userId"`
	} `dynamodbav:"notes"`
	CreatedAt int    `dynamodbav:"created_at"`
	AmountDue int    `dynamodbav:"amount_due"`
	Currency  string `dynamodbav:"currency"`
	Receipt   string `dynamodbav:"receipt"`
	ID        string `dynamodbav:"id"`
	OfferID   string `dynamodbav:"offer_id"`
	Entity    string `dynamodbav:"entity"`
	Attempts  int    `dynamodbav:"attempts"`
	Status    string `dynamodbav:"status"`
}

type RzpPayment struct {
	ID             string `dynamodbav:"id"`
	Entity         string `dynamodbav:"entity"`
	Amount         int    `dynamodbav:"amount"`
	Currency       string `dynamodbav:"currency"`
	Status         string `dynamodbav:"status"`
	OrderID        string `dynamodbav:"order_id"`
	InvoiceID      any    `dynamodbav:"invoice_id"`
	International  bool   `dynamodbav:"international"`
	Method         string `dynamodbav:"method"`
	AmountRefunded int    `dynamodbav:"amount_refunded"`
	RefundStatus   any    `dynamodbav:"refund_status"`
	Captured       bool   `dynamodbav:"captured"`
	Description    string `dynamodbav:"description"`
	CardID         any    `dynamodbav:"card_id"`
	Bank           string `dynamodbav:"bank"`
	Wallet         any    `dynamodbav:"wallet"`
	Vpa            any    `dynamodbav:"vpa"`
	Email          string `dynamodbav:"email"`
	Contact        string `dynamodbav:"contact"`
	Notes          struct {
		Mobile  string `dynamodbav:"mobile"`
		UserID  string `dynamodbav:"userId"`
		Address string `dynamodbav:"address"`
	} `dynamodbav:"notes"`
	Fee              int `dynamodbav:"fee"`
	Tax              int `dynamodbav:"tax"`
	ErrorCode        any `dynamodbav:"error_code"`
	ErrorDescription any `dynamodbav:"error_description"`
	ErrorSource      any `dynamodbav:"error_source"`
	ErrorStep        any `dynamodbav:"error_step"`
	ErrorReason      any `dynamodbav:"error_reason"`
	AcquirerData     struct {
		BankTransactionID string `dynamodbav:"bank_transaction_id"`
	} `dynamodbav:"acquirer_data"`
	CreatedAt int `dynamodbav:"created_at"`
}

type RzpRefunds struct {
	AcquirerData struct {
		Arn string `dynamodbav:"arn"`
	} `dynamodbav:"acquirer_data"`
	Amount    int    `dynamodbav:"amount"`
	BatchID   any    `dynamodbav:"batch_id"`
	CreatedAt int    `dynamodbav:"created_at"`
	Currency  string `dynamodbav:"currency"`
	Entity    string `dynamodbav:"entity"`
	ID        string `dynamodbav:"id"`
	Notes     struct {
		Address string `dynamodbav:"address"`
		Mobile  string `dynamodbav:"mobile"`
		UserID  string `dynamodbav:"userId"`
	} `dynamodbav:"notes"`
	PaymentID      string `dynamodbav:"payment_id"`
	Receipt        any    `dynamodbav:"receipt"`
	SpeedProcessed string `dynamodbav:"speed_processed"`
	SpeedRequested string `dynamodbav:"speed_requested"`
	Status         string `dynamodbav:"status"`
}
