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
	s3FileID := uuid.New().String()

	// Add entry in group table
	columnNames, err := getCsvColumnNames(file)
	if err != nil {
		log.Fatalf("error fetching column item: %v", err)
	}

	entryGroupItem := map[string]types.AttributeValue{
		"ID":   &types.AttributeValueMemberS{Value: s3FileID},
		"Name": &types.AttributeValueMemberS{Value: request.Name},

		// TODO fetch user id from token
		"UserID": &types.AttributeValueMemberS{Value: "USERID"},

		// Array of column names of csv file
		"ColumnNames": &types.AttributeValueMemberSS{Value: columnNames},
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
