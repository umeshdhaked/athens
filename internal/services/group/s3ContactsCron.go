package group

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/pkg/aws"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/pkg/cron"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/fastbiztech/hastinapura/pkg/logger"
	"github.com/fastbiztech/hastinapura/pkg/mutex"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	S3ContactsFetchProcessingBatchSize = 5

	// mutex lock keys
	MutexKeyS3ContactsFetchProcessing = "mutex-key-s3-contacts-fetch-processing-%s"
)

type s3ContactsCronExecutor struct {
	ctx                *gin.Context
	pendingJobsRepo    *repo.PendingJobsRepo
	cronProcessingRepo *repo.CronProcessingRepo
	contactsRepo       *repo.ContactsRepo
}

// CsvData represents the CSV data containing both column names and rows
type s3ContactsCsvData struct {
	ColumnNames []string
	Rows        [][]string
}

func InitialiseS3ContactsCron() {
	// Worker check
	// todo add worker args
	if os.Getenv(constants.WorkerCronArg) != constants.WorkerCronArgS3Contacts {
		return
	}

	if !config.GetConfig().Crons.CronsConfigS3Contacts.Enable {
		return
	}

	newCtx, _ := gin.CreateTestContext(nil)

	job := (&cron.Scheduler{}).NewScheduler()

	job.Initialize(
		time.Duration(config.GetConfig().Crons.CronsConfigS3Contacts.ExecutionTime)*time.Second,
		time.Duration(config.GetConfig().Crons.CronsConfigS3Contacts.StartTime)*time.Second,
		&s3ContactsCronExecutor{
			ctx:                newCtx,
			pendingJobsRepo:    repo.GetPendingJobsRepo(),
			cronProcessingRepo: repo.GetCronProcessingRepo(),
			contactsRepo:       repo.GetContactsRepo(),
		})
}

func (s *s3ContactsCronExecutor) JobExecutor() {
	var (
		toProcessPendingJob models.PendingJobs
		groupName           string
	)

	logger.GetLogger().
		WithField("type", "s3 contacts fetch").
		Info("cron run started")

	// fetch pending jobs
	pendingJobItems, err := s.pendingJobsRepo.FetchAllByConditions(s.ctx, dtos.GetPendingJobsRequest{
		Status: models.PendingJobsStatusPending,
		Type:   models.PendingJobsTypeSmsContactsPullFromS3,
	})
	if err != nil {
		fmt.Println("failed to fetch pending jobs in cron")
		return
	}

	if len(pendingJobItems) == 0 {
		fmt.Println("no pending jobs to run in cron")
		return
	}

	toProcessPendingJob = pendingJobItems[0]

	// fetch group name info from pending job entity
	if _, ok := toProcessPendingJob.Extra[constants.ConstGroupName].(string); !ok {
		fmt.Println("group name not attached to pending job of contacts creation")
		return
	}

	groupName = toProcessPendingJob.Extra[constants.ConstGroupName].(string)

	cronProcessingEntity, err := s.getCronProcessingRowNumber(s.ctx, toProcessPendingJob.Name, int(S3ContactsFetchProcessingBatchSize))
	if err != nil {
		fmt.Println("failed to fetch next row number to update")
		return
	}

	// get all the rows to be updated and check if they already exists in the contacts table.
	rowsToProcess, err := s.getRowsFromS3CsvFile(s.ctx, toProcessPendingJob.Name, cronProcessingEntity.InProgress+1, cronProcessingEntity.InProgress+int(S3ContactsFetchProcessingBatchSize))
	if err != nil {
		log.Println("failed fetching new rows from csv file")
		return
	}

	var toUpdateContacts []models.Contacts
	for _, rowToProcess := range rowsToProcess.Rows {
		var toUpdateContact models.Contacts

		toUpdateContact.ID = uuid.New().String()
		toUpdateContact.GroupName = groupName
		toUpdateContact.Additional = make(map[string]interface{})

		for indexColumn, column := range rowsToProcess.ColumnNames {
			switch strings.ToLower(column) {
			case models.ColumnContactsName:
				toUpdateContact.Name = rowToProcess[indexColumn]
			case models.ColumnContactsEmail:
				toUpdateContact.Email = rowToProcess[indexColumn]
			case models.ColumnContactsMobile:
				toUpdateContact.Mobile = rowToProcess[indexColumn]
			default:
				toUpdateContact.Additional[column] = rowToProcess[indexColumn]
			}
		}

		// check name/email/mobile combination if exists already in the database.
		dedupEntities, err := s.contactsRepo.FetchAllByConditions(
			s.ctx,
			dtos.GetContactsRequest{
				Name:   toUpdateContact.Name,
				Email:  toUpdateContact.Email,
				Mobile: toUpdateContact.Mobile,
			}, "")
		if err != nil || len(dedupEntities) > 0 {
			fmt.Println("either err fetching contacts or entry already exists")
			continue
		}

		toUpdateContacts = append(toUpdateContacts, toUpdateContact)
	}

	// create contacts
	err = s.contactsRepo.BulkCreate(s.ctx, toUpdateContacts)
	if err != nil {
		fmt.Println("failed contacts bulk insertion")
		return
	}

	// update cronProcessing item
	_, err = s.cronProcessingRepo.UpdateStatusByID(s.ctx, cronProcessingEntity.ID, models.CronProcessingStatusCompleted)
	if err != nil {
		fmt.Println("error in updating s3 processing entity")
		return
	}

	// if lesser rows to process then just update cronProcessing table and pendingJobs table and mark it completed
	if len(rowsToProcess.Rows) < S3ContactsFetchProcessingBatchSize {
		// update cronProcessing item
		_, err = s.cronProcessingRepo.UpdateStatusByID(s.ctx, cronProcessingEntity.ID, models.CronProcessingStatusLastRun)
		if err != nil {
			fmt.Println("error in updating s3 processing entity")
			return
		}

		// update pendingJobs item
		_, err = s.pendingJobsRepo.UpdateStatusByName(s.ctx, toProcessPendingJob.Name, models.PendingJobsStatusCompleted)
		if err != nil {
			fmt.Println("error in updating s3 processing entity")
			return
		}

		fmt.Printf("pending job completed: %s", toProcessPendingJob.Name)

		return
	}

	fmt.Printf("sucessfully created contacts. Processed row number: %v\n", cronProcessingEntity.InProgress)
}

func (s *s3ContactsCronExecutor) getCronProcessingRowNumber(ctx *gin.Context, s3FileName string, batch int) (models.CronProcessing, error) {

	// acquire lock on file name
	output, err := mutex.GetCronProcessingMutexLockManager().
		AcquireAndRelease(ctx,
			fmt.Sprintf(MutexKeyS3ContactsFetchProcessing, s3FileName),
			[]byte("Dummy Data"),
			func() (interface{}, error) {

				// get the last in progress row number
				cronProcessingItems, err := repo.GetCronProcessingRepo().FetchByName(ctx, s3FileName)
				if err != nil {
					fmt.Println("error in item fetch")
					return nil, err
				}

				// fetch last inProgress row
				lastInProgressRow := fetchLastInProgressRowForS3ContactsProcessing(cronProcessingItems)

				// insert entry for current batch
				cronProcessingInsertItem := models.CronProcessing{
					ID:         uuid.New().String(),
					Name:       s3FileName,
					Batch:      batch,
					InProgress: lastInProgressRow,
					Status:     models.CronProcessingStatusProcessing,
				}

				err = s.cronProcessingRepo.Create(ctx, &cronProcessingInsertItem)
				if err != nil {
					fmt.Println("error in cron processing item insertion")
					return models.CronProcessing{}, err
				}

				// return batch number
				return cronProcessingInsertItem, nil
			})

	if err != nil {
		// todo: mutex already acquire handling
		return models.CronProcessing{}, err
	}

	return output.(models.CronProcessing), nil
}

func (s *s3ContactsCronExecutor) getRowsFromS3CsvFile(ctx *gin.Context, s3FileName string, from, to int) (*s3ContactsCsvData, error) {
	var (
		body io.Reader
		err  error
	)

	// pull s3 file
	body, err = aws.GetS3Client().Fetch(ctx, aws.BucketContactUpload, s3FileName)
	if err != nil {
		log.Println("error fetching s3 file")
		return nil, err
	}

	// Parse the CSV data from the response body
	reader := csv.NewReader(body)

	// Extract column names from the first row
	columnNames, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %v", err)
	}

	// Extract rows within the specified range
	var extractedRows [][]string
	rowNum := 1 // since we have already read the column names
	for rowNum <= to {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %v", err)
		}

		// Check if the current row number is within the specified range
		if rowNum >= from && rowNum <= to {
			extractedRows = append(extractedRows, record)
		}

		// Increment the row number
		rowNum++
	}

	return &s3ContactsCsvData{
		ColumnNames: columnNames,
		Rows:        extractedRows,
	}, nil
}

func fetchLastInProgressRowForS3ContactsProcessing(models []models.CronProcessing) int {
	var lastInProgressRow int

	if len(models) == 0 {
		return 0
	}

	for _, model := range models {
		if lastInProgressRow < (model.InProgress + model.Batch) {
			lastInProgressRow = model.InProgress + model.Batch
		}
	}

	return lastInProgressRow
}
