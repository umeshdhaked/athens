package group

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"

	"github.com/umeshdhaked/athens/internal/constants"
	"github.com/umeshdhaked/athens/internal/models"
	"github.com/umeshdhaked/athens/internal/pkg/aws"
	"github.com/umeshdhaked/athens/internal/pkg/repo"
	"github.com/umeshdhaked/athens/internal/utils"
	"github.com/umeshdhaked/athens/pkg/dtos"
	"github.com/umeshdhaked/athens/pkg/logger"
	"github.com/gin-gonic/gin"
	gormLogger "gorm.io/gorm/logger"
)

type Service struct {
	baseRepo           repo.IRepository
	groupRepo          repo.IGroupRepo
	cronProcessingRepo repo.ICronProcessingRepo
	pendingJobsRepo    repo.IPendingJobsRepo
}

var (
	once                     sync.Once
	service                  *Service
	validGroupContactsColumn = []string{
		constants.Name,
		constants.Mobile,
		constants.Email,
	}
)

func InitialiseService() {
	once.Do(func() {
		service = &Service{
			baseRepo:           repo.GetRepository(),
			groupRepo:          repo.GetGroupRepo(),
			cronProcessingRepo: repo.GetCronProcessingRepo(),
			pendingJobsRepo:    repo.GetPendingJobsRepo(),
		}
	})
}

func GetService() *Service {
	return service
}

func (s *Service) UploadGroupToS3(c *gin.Context, file multipart.File, request dtos.UploadGroupContactsRequest) (interface{}, error) {
	//  add validation for same name already existed
	var items []models.ContactGroups
	err := s.baseRepo.FindMultiple(c, &items, map[string]interface{}{
		constants.UserID: c.GetInt64(constants.JwtTokenUserID),
		constants.Name:   request.Name,
	})
	if err != nil && !errors.Is(err, gormLogger.ErrRecordNotFound) {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	if len(items) > 0 {
		logger.GetLogger().Error("duplicate UserID/Name entry")
		return nil, errors.New("duplicate UserID/Name entry")
	}

	// Add entry in group table
	columnNames, err := getCsvColumnNames(file)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	if !isValidColumnNamesForGroupContacts(columnNames) {
		logger.GetLogger().Error("invalid columns")
		return nil, errors.New("invalid columns")
	}

	columnNamesBytes, _ := json.Marshal(columnNames)
	// Insert item into the database
	entryGroupItem := models.ContactGroups{
		Name:        request.Name,
		UserID:      c.GetInt64(constants.JwtTokenUserID),
		ColumnNames: string(columnNamesBytes),
	}

	err = s.baseRepo.Create(c, &entryGroupItem)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	// Upload file to S3
	if err := aws.GetS3Client().Upload(file, aws.BucketContactUpload, strconv.FormatInt(entryGroupItem.ID, 10)); err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	// add pending job entry
	extraData, _ := json.Marshal(map[string]interface{}{
		constants.ConstGroupName: request.Name,
	})
	entryPendingJobsItem := models.PendingJobs{
		Name:   strconv.FormatInt(entryGroupItem.ID, 10),
		Type:   models.PendingJobsTypeSmsContactsPullFromS3,
		Status: models.PendingJobsStatusPending,
		Extra:  extraData,
	}

	err = s.baseRepo.Create(c, &entryPendingJobsItem)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	return nil, nil
}

func getCsvColumnNames(file multipart.File) ([]string, error) {
	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the first record which contains column names
	columnNames, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV record: %v", err)
	}

	// Reset file pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		// Handle error
	}

	return columnNames, nil
}

func isValidColumnNamesForGroupContacts(columnNames []string) bool {

	for _, column := range validGroupContactsColumn {
		if !utils.Contains(columnNames, strings.Title(column)) {
			return false
		}
	}

	return true
}

func (s *Service) GetContacts(c *gin.Context, request dtos.GetGroupRequest) (interface{}, error) {
	var items []models.ContactGroups
	err := s.baseRepo.FindMultiplePagination(c, &items, map[string]interface{}{
		constants.Name:   request.Name,
		constants.UserID: request.UserID,
	}, dtos.Pagination{
		From: request.From,
		To:   request.To,
	})
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("error fetching item:")
	}

	return items, nil
}
