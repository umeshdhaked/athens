package models

const (
	TablePendingJobs = "pendingJobs"
)

const (
	ColumnPendingJobsName   = "name"
	ColumnPendingJobsType   = "type"
	ColumnPendingJobsStatus = "status"

	// Types
	PendingJobsTypeSmsContactsPullFromS3 = "sms_contacts_pull_from_s3"

	// Statuses
	PendingJobsStatusPending   = "PENDING"
	PendingJobsStatusCompleted = "COMPLETED"
)

type PendingJobs struct {
	Name   string                 `json:"name" dynamodbav:"name"`
	Type   string                 `json:"type" dynamodbav:"type"`
	Status string                 `json:"status" dynamodbav:"status"`
	Extra  map[string]interface{} `json:"extra" dynamodbav:"extra"`
	BaseModel
}
