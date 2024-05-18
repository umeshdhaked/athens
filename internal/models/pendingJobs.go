package models

import (
	"gorm.io/datatypes"
)

const (
	TablePendingJobs = "pending_jobs"
)

const (

	// Types
	PendingJobsTypeSmsContactsPullFromS3 = "sms_contacts_pull_from_s3"

	// Statuses
	PendingJobsStatusPending   = "PENDING"
	PendingJobsStatusCompleted = "COMPLETED"
)

type PendingJobs struct {
	ID     int64          `json:"id" dynamodbav:"id"`
	Name   string         `json:"name" dynamodbav:"name"`
	Type   string         `json:"type" dynamodbav:"type"`
	Status string         `json:"status" dynamodbav:"status"`
	Extra  datatypes.JSON `json:"extra" dynamodbav:"extra"`
	BaseModel
}

func (o *PendingJobs) TableName() string {
	return TablePendingJobs
}

func (o *PendingJobs) GetID() int64 {
	return o.ID
}

func (o *PendingJobs) SetID(id int64) {
	o.ID = id
}
