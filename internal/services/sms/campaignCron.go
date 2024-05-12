package sms

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fastbiztech/hastinapura/internal/config"
	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/pkg/cron"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/fastbiztech/hastinapura/pkg/logger"
	"github.com/fastbiztech/hastinapura/pkg/mutex"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	campaignCronBatchSize = 5

	CronProcessingSmsCampaignPrefix = "sms_campaign_run"

	// mutex lock keys
	MutexKeyCampaignProcessing = "mutex-key-campaign-processing-%s"
)

type campaignCronExecutor struct {
	ctx                *gin.Context
	creditsRepo        *repo.CreditsRepo
	smsAuditRepo       *repo.SmsAuditRepo
	contactsRepo       *repo.ContactsRepo
	smsCampaignRepo    *repo.SmsCampaignRepo
	pendingJobsRepo    *repo.PendingJobsRepo
	smsTemplateRepo    *repo.SmsTemplateRepo
	cronProcessingRepo *repo.CronProcessingRepo
}

func InitialiseCampaignCron() {
	// Worker check
	// todo add worker args
	if os.Getenv(constants.WorkerCronArg) != constants.WorkerCronArgCampaign {
		return
	}

	if !config.GetConfig().Crons.CronsConfigCampaign.Enable {
		return
	}

	newCtx, _ := gin.CreateTestContext(nil)

	job := (&cron.Scheduler{}).NewScheduler()

	job.Initialize(
		time.Duration(config.GetConfig().Crons.CronsConfigCampaign.ExecutionTime)*time.Second,
		time.Duration(config.GetConfig().Crons.CronsConfigCampaign.StartTime)*time.Second,
		&campaignCronExecutor{
			ctx:                newCtx,
			creditsRepo:        repo.GetCreditsRepo(),
			smsAuditRepo:       repo.GetSmsAuditRepo(),
			contactsRepo:       repo.GetContactsRepo(),
			smsCampaignRepo:    repo.GetSmsCampaignRepo(),
			pendingJobsRepo:    repo.GetPendingJobsRepo(),
			smsTemplateRepo:    repo.GetSmsTemplateRepo(),
			cronProcessingRepo: repo.GetCronProcessingRepo(),
		})
}

func (s *campaignCronExecutor) JobExecutor() {

	logger.GetLogger().
		WithField("type", "campaign processing").
		Info("cron run started")

	// TODO Dont start campaign until all contacts all loaded from s3

	// Check which scheduled campaign to run
	campaigns, err := getCampaignsToRun(s.ctx)
	if err != nil {
		logger.GetLogger().
			WithField("error", err).
			Info("failed to fetch campaigns to run")
		return
	}

	var wg sync.WaitGroup // Create a WaitGroup
	for _, campaign := range campaigns {
		wg.Add(1)

		go func(smsCampaign models.SmsCampaign) {
			defer wg.Done()

			// Run campaign only if s3 contacts upload is successful
			smsTemplate, err := repo.GetSmsTemplateRepo().FetchByIDAndUserID(s.ctx, smsCampaign.TemplateID, smsCampaign.UserID)
			if err != nil {
				logger.GetLogger().Error(err.Error())
				return
			}

			if smsTemplate.Status != models.SmsTemplateStateApproved {
				logger.GetLogger().Error("sms template not approved")
				return
			}

			// Check credits left
			creditEntity, creditCost, allowed := isEnoughBalanceForSmsCampaignRun(s.ctx, smsCampaign.UserID, smsCampaign.TemplateID,
				constants.SubCategoryPromotional,
				campaignCronBatchSize)
			if !allowed {
				logger.GetLogger().Error("balance check failed for sms campaign run")
				return
			}

			// Fetch which contacts to pick for message processing
			cronProcessingEntity, err := s.getCronProcessingRowNumber(s.ctx, smsCampaign.ID, campaignCronBatchSize)
			if err != nil {
				logger.GetLogger().Error("failed to fetch next row number to update")
				return
			}

			contactEntities, err := s.getContactsForCronProcessing(s.ctx, smsCampaign, cronProcessingEntity.LastEvaluatedID,
				campaignCronBatchSize)
			if err != nil {
				return
			}

			// mark campaign completed if it is
			if len(contactEntities) == 0 {
				err = s.smsCampaignRepo.UpdateByID(s.ctx, smsCampaign.ID, models.TableSmsCampaign, map[string]interface{}{
					models.ColumnSmsCampaignStatus: models.SmsCampaignStateExecuted,
				}, &smsCampaign)
				if err != nil {
					logger.GetLogger().Error(err.Error())
					return
				}
			}

			// TODO: hit for actual sms send

			// Mark completed contacts
			// And save lastEvaluated key in cron processing
			err = s.cronProcessingRepo.UpdateByID(s.ctx, cronProcessingEntity.ID, models.TableCronProcessing,
				map[string]interface{}{
					models.ColumnCronProcessingLastEvaluatedID: contactEntities[len(contactEntities)-1].ID,
					models.ColumnCronProcessingStatus:          models.CronProcessingStatusCompleted,
				}, &cronProcessingEntity)
			if err != nil {
				logger.GetLogger().Error("failed to update cron processing entity after sms processing")
				return
			}

			// Consume credits
			_, err = s.creditsRepo.UpdateCreditsLeftByID(s.ctx, creditEntity.ID, creditEntity.CreditsLeft-creditCost)
			if err != nil {
				return
			}

			// Audit sms sent to contacts
			var bulkCreateSmsAuditEntities []models.SmsAudit
			for _, contactEntity := range contactEntities {

				smsAudit := models.SmsAudit{
					ID:              uuid.New().String(),
					UserID:          smsCampaign.UserID,
					CreditsConsumed: creditCost / float64(campaignCronBatchSize),
					TemplateID:      smsCampaign.TemplateID,
					SenderCode:      smsCampaign.SenderCode,
					ContactID:       contactEntity.ID,
					Status:          "PROCESSED",
					TriggeredMode:   "cron",
				}

				bulkCreateSmsAuditEntities = append(bulkCreateSmsAuditEntities, smsAudit)
			}

			err = s.smsAuditRepo.BulkCreate(s.ctx, bulkCreateSmsAuditEntities)
			if err != nil {
				return
			}

		}(campaign)
	}

	wg.Wait()

	logger.GetLogger().
		WithField("type", "campaign processing").
		Info("cron run completed")
}

func (s *campaignCronExecutor) getContactsForCronProcessing(ctx *gin.Context,
	smsCampaign models.SmsCampaign,
	lastEvaluatedID string,
	batch int) ([]models.Contacts, error) {

	contactEntities, err := s.contactsRepo.FetchAllByConditions(ctx, dtos.GetContactsRequest{
		GroupName: smsCampaign.GroupName,
		Limit:     batch,
	}, lastEvaluatedID)
	if err != nil {
		logger.GetLogger().
			Error("failed to fetch contacts")
		return nil, err
	}

	return contactEntities, nil
}

func (s *campaignCronExecutor) getCronProcessingRowNumber(ctx *gin.Context, smsCampaignID string, batch int) (models.CronProcessing, error) {
	// acquire lock on file name
	output, err := mutex.GetCronProcessingMutexLockManager().
		AcquireAndRelease(ctx,
			fmt.Sprintf(MutexKeyCampaignProcessing, smsCampaignID),
			[]byte("Dummy Data"),
			func() (interface{}, error) {

				// get the last in progress row number
				cronProcessingItems, err := repo.GetCronProcessingRepo().FetchByName(ctx,
					CronProcessingSmsCampaignPrefix+"_"+smsCampaignID)
				if err != nil && err.Error() != repo.ErrCodeNoDataFound {
					fmt.Println("error in item fetch")
					return nil, err
				}

				// fetch last inProgress row
				lastInProgressRow := fetchLastInProgressRowForSmsCampaignProcessing(cronProcessingItems)

				// insert entry for current batch
				cronProcessingInsertItem := models.CronProcessing{
					ID:         uuid.New().String(),
					Name:       CronProcessingSmsCampaignPrefix + "_" + smsCampaignID,
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

func fetchLastInProgressRowForSmsCampaignProcessing(models []models.CronProcessing) int {
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
