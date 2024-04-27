package models

const (
	TableS3Processing = "s3Processing"
)

const (
	ColumnS3ProcessingId         = "id"
	ColumnS3ProcessingName       = "name"
	ColumnS3ProcessingBatch      = "batch"
	ColumnS3ProcessingInProgress = "in_progress"
	ColumnS3ProcessingStatus     = "status"

	// statuses
	S3ProcessingStatusProcessing = "PROCESSING"
	S3ProcessingStatusCompleted  = "COMPLETED"
	S3ProcessingStatusLastRun    = "LAST_RUN"
)

type S3Processing struct {
	ID         string `json:"id" dynamodbav:"id"`
	Name       string `json:"name" dynamodbav:"name"`
	Batch      int    `json:"batch" dynamodbav:"batch"`
	InProgress int    `json:"in_progress" dynamodbav:"in_progress"`
	Status     string `json:"status" dynamodbav:"status"`
	BaseModel
}
