package models

const (
	// Base columns
	ColumnCreatedAt = "CreatedAt"
	ColumnUpdatedAt = "UpdatedAt"
	ColumnDeletedAt = "DeletedAt"

	// Add all index names here
	// format: index_<tableName>_<columnName>

	IndexTableGroupIndexName       = "index_group_name"
	IndexTableSmsSenderIndexUserID = "index_smssender_userid"

	// indexes: sms template table
	IndexTableSmsTemplateIndexUserID = "index_smstemplate_userid"

	// indexes: credits table
	IndexTableCreditsIndexUserID = "index_credits_userid"

	// indexes: user subscription table
	IndexTableUserSubscriptionIndexUserID = "index_usersub_userid"

	// indexes: sms campaign table
	IndexTableSmsCampaignIndexUserID = "index_smscampaign_userid"
	IndexTableSmsCampaignIndexStatus = "index_smscampaign_status"
)

type BaseModel struct {
	CreatedAt int64 `json:"created_at"` // Unix timestamps
	UpdatedAt int64 `json:"updated_at"`
	DeletedAt int64 `json:"deleted_at"`
}
