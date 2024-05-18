package sms

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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
	gormLogger "gorm.io/gorm/logger"
)

const (
	campaignCronBatchSize = 5

	CronProcessingSmsCampaignPrefix = "sms_campaign_run"

	// mutex lock keys
	MutexKeyCampaignProcessing    = "mutex-key-campaign-processing-%s"
	MutexKeyCampaignProcessingTTL = 30
)

type campaignCronExecutor struct {
	ctx                *gin.Context
	baseRepo           repo.IRepository
	creditsRepo        repo.ICreditsRepo
	smsAuditRepo       repo.ISmsAuditRepo
	contactsRepo       repo.IContactsRepo
	smsCampaignRepo    repo.ISmsCampaignRepo
	pendingJobsRepo    repo.IPendingJobsRepo
	smsTemplateRepo    repo.ISmsTemplateRepo
	cronProcessingRepo repo.ICronProcessingRepo
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
			baseRepo:           repo.GetRepository(),
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
			var smsTemplate models.SmsTemplate
			err := repo.GetRepository().Find(s.ctx, &smsCampaign, map[string]interface{}{
				constants.ID:     smsCampaign.TemplateID,
				constants.UserID: smsCampaign.UserID,
			})
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
				smsCampaign.Status = models.SmsCampaignStateExecuted
				err = s.baseRepo.Update(s.ctx, &smsCampaign)
				if err != nil {
					logger.GetLogger().Error(err.Error())
					return
				}
			}

			// TODO: hit for actual sms send

			// Mark completed contacts
			// And save lastEvaluated key in cron processing
			cronProcessingEntity.LastEvaluatedID = contactEntities[len(contactEntities)-1].ID
			cronProcessingEntity.Status = models.CronProcessingStatusCompleted
			err = s.baseRepo.Update(s.ctx, &cronProcessingEntity)
			if err != nil {
				logger.GetLogger().Error("failed to update cron processing entity after sms processing")
				return
			}

			// Consume credits
			creditEntity.BalanceLeft = creditEntity.BalanceLeft - creditCost
			err = s.baseRepo.Update(s.ctx, creditEntity)
			if err != nil {
				logger.GetLogger().Error(err.Error())
				return
			}

			// Audit sms sent to contacts
			var bulkCreateSmsAuditEntities []models.SmsAudit
			for _, contactEntity := range contactEntities {

				smsAudit := models.SmsAudit{
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
	lastEvaluatedID int64,
	batch int) ([]models.Contacts, error) {

	var contactEntities []models.Contacts
	// todo: build pagination thing
	err := s.baseRepo.FindMultiplePagination(ctx, &contactEntities, map[string]interface{}{
		constants.Name: smsCampaign.GroupName,
	}, dtos.Pagination{})
	if err != nil {
		logger.GetLogger().Error("failed to fetch contacts")
		return nil, err
	}

	return contactEntities, nil
}

func (s *campaignCronExecutor) getCronProcessingRowNumber(ctx *gin.Context, smsCampaignID int64, batch int) (models.CronProcessing, error) {
	// acquire lock on file name
	output, err := mutex.GetClient().AcquireAndRelease(ctx,
		fmt.Sprintf(MutexKeyCampaignProcessing, smsCampaignID),
		MutexKeyCampaignProcessingTTL*time.Second,
		func() (interface{}, error) {

			// get the last in progress row number
			var cronProcessingItems []models.CronProcessing
			err := s.baseRepo.FindMultiple(ctx, &cronProcessingItems, map[string]interface{}{
				constants.Name: CronProcessingSmsCampaignPrefix + "_" + strconv.FormatInt(smsCampaignID, 10),
			})
			if err != nil && !errors.Is(err, gormLogger.ErrRecordNotFound) {
				fmt.Println("error in item fetch")
				return nil, err
			}

			// fetch last inProgress row
			lastInProgressRow := fetchLastInProgressRowForSmsCampaignProcessing(cronProcessingItems)

			// insert entry for current batch
			cronProcessingInsertItem := models.CronProcessing{
				Name:       CronProcessingSmsCampaignPrefix + "_" + strconv.FormatInt(smsCampaignID, 10),
				Batch:      batch,
				InProgress: lastInProgressRow,
				Status:     models.CronProcessingStatusProcessing,
			}

			err = s.baseRepo.Create(ctx, &cronProcessingInsertItem)
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
