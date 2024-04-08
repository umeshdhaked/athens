package models

const (
	// Add all index names here
	// format: index_<tableName>_<columnName>

	IndexTableGroupIndexName       = "index_group_name"
	IndexTableSmsSenderIndexUserID = "index_smssender_userid"

	// indexes: sms template table
	IndexTableSmsTemplateIndexUserID = "index_smstemplate_userid"

	// indexes: credits table
	IndexTableCreditsIndexUserID = "index_credits_userid"

	// indexes: credits table
	IndexTableUserSubscriptionIndexUserID = "index_usersub_userid"
)

type BaseModel struct {
	CreatedAt int `json:"created_at"` // Unix timestamps
	UpdatedAt int `json:"updated_at"`
	DeletedAt int `json:"deleted_at"`
}
