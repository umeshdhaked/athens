package repo

import (
	"context"
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/pkg/db"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/fastbiztech/hastinapura/pkg/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const (
	// Error codes
	ErrCodeNoDataFound = "no_data_found"
)

var (
	once     sync.Once
	baseRepo IRepository
)

type IRepository interface {
	Create(ctx context.Context, model models.IModel) error
	CreateInBatches(ctx context.Context, models interface{}, batch int) error
	Update(ctx *gin.Context, model models.IModel) error
	Delete(ctx *gin.Context, model models.IModel) error
	FindByID(ctx *gin.Context, id interface{}, model models.IModel) error
	Find(ctx *gin.Context, model models.IModel, condition map[string]interface{}) error
	FindMultiple(ctx *gin.Context, models interface{}, condition map[string]interface{}) error
	FindMultiplePagination(ctx *gin.Context,
		models interface{},
		condition map[string]interface{},
		pagination dtos.Pagination) error
}

type Repository struct {
	mysqlDB  *gorm.DB
	dynamoDB *dynamodb.Client
}

func GetRepository() IRepository {
	return baseRepo
}

func newRepository(mysqlDbClient *gorm.DB, dynamoDbClient *dynamodb.Client) {
	baseRepo = &MysqlRepository{
		db: mysqlDbClient,
	}
}

func InitialiseRepositories() {

	mysqlDBClient := db.GetDb().Mysql

	once.Do(func() {
		newRepository(mysqlDBClient, nil)
		newCreditsAuditRepo(mysqlDBClient, nil)
		newCreditsRepo(mysqlDBClient, nil)
		newGroupRepo(mysqlDBClient, nil)
		newPendingJobsRepo(mysqlDBClient, nil)
		newCronProcessingRepo(mysqlDBClient, nil)
		newContactsRepo(mysqlDBClient, nil)
		newOtpRepo(mysqlDBClient, nil)
		newPricingRepo(mysqlDBClient, nil)
		newPromotionRepo(mysqlDBClient, nil)
		newSmsAuditRepo(mysqlDBClient, nil)
		newSmsCampaignRepo(mysqlDBClient, nil)
		newSmsSenderRepo(mysqlDBClient, nil)
		newSmsTemplateRepo(mysqlDBClient, nil)
		newSubscriptionRepo(mysqlDBClient, nil)
		newSubscriptionRepo(mysqlDBClient, nil)
		newUserRepo(mysqlDBClient, nil)
		newKycRepo(mysqlDBClient, nil)
		newPaymentsRepo(mysqlDBClient, nil)
		newInvoiceRepo(mysqlDBClient, nil)
	})
}

func GetError(ctx context.Context, query *gorm.DB) error {

	if errors.Is(query.Error, gormLogger.ErrRecordNotFound) {
		logger.GetLogger().Info(ErrCodeNoDataFound)
		return query.Error
	}

	// extract custom error message & handle
	if query.Error != nil {
		logger.GetLogger().Info(query.Error.Error())
	}

	return query.Error
}
