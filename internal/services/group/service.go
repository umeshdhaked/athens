package group

import (
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"sync"

	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/pkg/aws"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/internal/utils"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service struct {
	baseRepo         *repo.Repository
	groupRepo        *repo.GroupRepo
	s3ProcessingRepo *repo.S3ProcessingRepo
	pendingJobsRepo  *repo.PendingJobsRepo
}

var (
	once                     sync.Once
	service                  *Service
	validGroupContactsColumn = []string{
		models.ColumnContactsName,
		models.ColumnContactsMobile,
		models.ColumnContactsEmail,
	}
)

func InitialiseService() {
	once.Do(func() {
		service = &Service{
			baseRepo:         repo.GetRepository(),
			groupRepo:        repo.GetGroupRepo(),
			s3ProcessingRepo: repo.GetS3ProcessingRepo(),
			pendingJobsRepo:  repo.GetPendingJobsRepo(),
		}
	})
}

func GetService() *Service {
	return service
}

func (s *Service) UploadGroupToS3(c *gin.Context, file multipart.File, request dtos.UploadGroupContactsRequest) (interface{}, error) {
	//  add validation for same name already existed
	items, err := s.groupRepo.FetchByUserIDAndName(c, c.GetString(constants.JwtTokenUserID), request.Name)
	if err != nil {
		log.Fatalf("error fetching column item: %v", err)
	}

	if len(items) > 0 {
		log.Fatalf("duplicate UserID/Name entry")
	}

	s3FileID := uuid.New().String()

	// Add entry in group table
	columnNames, err := getCsvColumnNames(file)
	if err != nil {
		log.Fatalf("error fetching column item: %v", err)
	}

	if !isValidColumnNamesForGroupContacts(columnNames) {
		log.Fatalf("invalid columns")
	}

	// Insert item into the database
	entryGroupItem := models.Group{
		ID:          s3FileID,
		Name:        request.Name,
		UserID:      c.GetString(constants.JwtTokenUserID),
		ColumnNames: columnNames,
	}

	err = s.groupRepo.CreateGroup(c, &entryGroupItem)
	if err != nil {
		log.Fatalf("error inserting item: %v", err)
	}

	// Upload file to S3
	if err := aws.GetS3Client().Upload(file, aws.BucketContactUpload, s3FileID); err != nil {
		return nil, err
	}

	// add pending job entry
	entryPendingJobsItem := models.PendingJobs{
		Name:   s3FileID,
		Type:   models.PendingJobsTypeSmsContactsPullFromS3,
		Status: models.PendingJobsStatusPending,
		Extra: map[string]interface{}{
			constants.ConstGroupName: request.Name,
		},
	}

	err = s.pendingJobsRepo.Create(c, &entryPendingJobsItem)
	if err != nil {
		log.Fatalf("error inserting item: %v", err)
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
	items, err := s.groupRepo.FetchAllByConditions(c, request)
	if err != nil {
		log.Fatalf("error fetching item: %v", err)
	}

	return items, nil
}
