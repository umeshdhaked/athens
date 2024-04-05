package group

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/pkg/aws"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/pkg/models/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service struct {
	baseRepo *repo.Repository
}

var (
	once    sync.Once
	service *Service
)

func InitialiseService() {
	once.Do(func() {
		service = &Service{
			baseRepo: repo.GetRepository(),
		}
	})
}

func (s *Service) UploadGroupToS3(c *gin.Context, file multipart.File, request dtos.UploadGroupContactsRequest) (interface{}, error) {
	//  add validation for same name already existed
	conditions := dtos.DbConditions{
		Index: "test_index1",
		PKey: map[string]interface{}{
			models.ColumnName: request.Name,
		},
		NonPKey: map[string]interface{}{
			models.ColumnUserID: "USERID",
		},
	}

	items, err := s.baseRepo.QueryItems(context.Background(), models.TableGroup, conditions)
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

	entryGroupItem := map[string]types.AttributeValue{
		models.ColumnID:   &types.AttributeValueMemberS{Value: s3FileID},
		models.ColumnName: &types.AttributeValueMemberS{Value: request.Name},

		// TODO fetch user id from token
		models.ColumnUserID: &types.AttributeValueMemberS{Value: "USERID"},

		// Array of column names of csv file
		models.ColumnColumnNames: &types.AttributeValueMemberSS{Value: columnNames},
	}

	// Insert item into the database
	err = s.baseRepo.CreateItem(context.Background(), models.TableGroup, entryGroupItem)
	if err != nil {
		log.Fatalf("error inserting item: %v", err)
	}

	// Upload file to S3
	if err := aws.GetS3Client().Upload(file, aws.BucketContactUpload, s3FileID); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Service) GetGroupContacts(c *gin.Context, request dtos.GetGroupContactsRequest) (interface{}, error) {
	conditions := dtos.DbConditions{
		Index: "test_index1",
		PKey: map[string]interface{}{
			models.ColumnName: request.Name,
		},
	}

	// Insert item into the database
	items, err := s.baseRepo.QueryItems(context.Background(), models.TableGroup, conditions)
	if err != nil {
		log.Fatalf("error inserting item: %v", err)
	}

	return items, nil
}

func GetService() *Service {
	return service
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
