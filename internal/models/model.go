package models

const (
	// Base columns
	ColumnCreatedAt = "CreatedAt"
	ColumnUpdatedAt = "UpdatedAt"
	ColumnDeletedAt = "DeletedAt"

	// Add all index names here
	// format: index_<tableName>_<columnName>

	// indexes: group table
	IndexTableGroupIndexUserID = "index_group_userid"

	// indexes: sms sender table
	IndexTableSmsSenderIndexUserID = "index_smssender_userid"

	// indexes: sms template table
	IndexTableSmsTemplateIndexUserID = "index_smstemplate_userid"

	// indexes: credits table
	IndexTableCreditsIndexUserID = "index_credits_userid"

	// indexes: subscription table
	IndexTableSubscriptionIndexUserID = "index_subscription_userid"

	// indexes: sms campaign table
	IndexTableSmsCampaignIndexUserID = "index_smscampaign_userid"
	IndexTableSmsCampaignIndexStatus = "index_smscampaign_status"

	//indexes: otp table
	IndexTableOtpIndexMobile = "index_otp_mobile"

	//indexes: pricing table
	IndexTablePricingIndexCategory     = "index_pricing_category"
	IndexTablePricingIndexPricingState = "index_pricing_pricingState"

	//indexes: user table
	IndexTableUserIndexMobile = "index_user_mobile"

	//indexes: promo-phones-no table
	IndexTablePromoPhoneIndexIsAlreadyContacted = "index_promo-phones-no_isAlreadyContacted"
	IndexTablePromoPhoneIndexMobile             = "index_promo-phones-no_mobile"

	// indexes: cronProcessing table
	IndexTableCronProcessingIndexName = "index_cron_processing_name"

	// indexes: RzpOrder table
	IndexTableRzpOrdersIndexStatus    = "index_rzpOrders_status"
	IndexTableRzpOrdersIndexCreatedAt = "index_rzpOrders_status_created_at"

	// indexes: RzpPayment table
	IndexTableRzpPaymentIndexOrderId = "index_rzpPayments_orderId"

	// indexes: RzpRefunds table
	IndexTableRzpRefundsIndexPaymentId = "index_rzpRefunds_paymentId"

	// indexes: FBTPayment table
	IndexTableFBTPaymentIndexRazorpayOrderId = "index_rzpPayments_razorpayOrderId"

	// indexes: Invoice table
	IndexTableInvoicesIndexInvoiceId = "index_invoices_invoiceId"
	IndexTableInvoicesIndexOrderId   = "index_invoices_orderId"
)

type BaseModel struct {
	CreatedAt int64 `gorm:"column:created_at" json:"created_at"` // Unix timestamps
	UpdatedAt int64 `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt int64 `gorm:"column:deleted_at" json:"deleted_at"`
}
