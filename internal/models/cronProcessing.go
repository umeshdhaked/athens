package models

const (
	TableCronProcessing = "cron_processing"
)

const (
	// statuses
	CronProcessingStatusProcessing = "PROCESSING"
	CronProcessingStatusCompleted  = "COMPLETED"
	CronProcessingStatusLastRun    = "LAST_RUN"
)

type CronProcessing struct {
	ID              int64  `json:"id" dynamodbav:"id"`
	Name            string `json:"name" dynamodbav:"name"`
	Batch           int    `json:"batch" dynamodbav:"batch"`
	InProgress      int    `json:"in_progress" dynamodbav:"in_progress"`
	Status          string `json:"status" dynamodbav:"status"`
	LastEvaluatedID int64  `json:"last_evaluated_id" dynamodbav:"last_evaluated_id"`
	BaseModel
}

func (o *CronProcessing) TableName() string {
	return TableCronProcessing
}

func (o *CronProcessing) GetID() int64 {
	return o.ID
}

func (o *CronProcessing) SetID(id int64) {
	o.ID = id
}
