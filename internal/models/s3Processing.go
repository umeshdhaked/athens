package models

const (
	TableCronProcessing = "cron_processing"
)

const (
	ColumnCronProcessingId              = "id"
	ColumnCronProcessingName            = "name"
	ColumnCronProcessingBatch           = "batch"
	ColumnCronProcessingInProgress      = "in_progress"
	ColumnCronProcessingStatus          = "status"
	ColumnCronProcessingLastEvaluatedID = "last_evaluated_id"

	// statuses
	CronProcessingStatusProcessing = "PROCESSING"
	CronProcessingStatusCompleted  = "COMPLETED"
	CronProcessingStatusLastRun    = "LAST_RUN"
)

type CronProcessing struct {
	ID              string `json:"id" dynamodbav:"id"`
	Name            string `json:"name" dynamodbav:"name"`
	Batch           int    `json:"batch" dynamodbav:"batch"`
	InProgress      int    `json:"in_progress" dynamodbav:"in_progress"`
	Status          string `json:"status" dynamodbav:"status"`
	LastEvaluatedID string `json:"last_evaluated_id" dynamodbav:"last_evaluated_id"`
	BaseModel
}
