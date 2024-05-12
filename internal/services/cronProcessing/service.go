package cronProcessing

import (
	"log"
	"sync"

	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

type Service struct {
	baseRepo           *repo.Repository
	cronProcessingRepo *repo.CronProcessingRepo
}

var (
	once    sync.Once
	service *Service
)

func InitialiseService() {
	once.Do(func() {
		service = &Service{
			baseRepo:           repo.GetRepository(),
			cronProcessingRepo: repo.GetCronProcessingRepo(),
		}
	})
}

func GetService() *Service {
	return service
}

func (s *Service) GetCronProcessing(c *gin.Context, request dtos.GetCronProcessingRequest) (interface{}, error) {
	items, err := s.cronProcessingRepo.FetchAllByConditions(c, request)
	if err != nil {
		log.Fatalf("error fetching item: %v", err)
	}

	return items, nil
}
