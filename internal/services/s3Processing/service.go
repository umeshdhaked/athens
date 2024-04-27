package s3Processing

import (
	"log"
	"sync"

	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

type Service struct {
	baseRepo         *repo.Repository
	s3ProcessingRepo *repo.S3ProcessingRepo
}

var (
	once    sync.Once
	service *Service
)

func InitialiseService() {
	once.Do(func() {
		service = &Service{
			baseRepo:         repo.GetRepository(),
			s3ProcessingRepo: repo.GetS3ProcessingRepo(),
		}
	})
}

func GetService() *Service {
	return service
}

func (s *Service) GetS3Processing(c *gin.Context, request dtos.GetS3ProcessingRequest) (interface{}, error) {
	items, err := s.s3ProcessingRepo.FetchAllByConditions(c, request)
	if err != nil {
		log.Fatalf("error fetching item: %v", err)
	}

	return items, nil
}
