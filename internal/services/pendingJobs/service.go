package pendingJobs

import (
	"log"
	"sync"

	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

type Service struct {
	baseRepo        *repo.Repository
	pendingJobsRepo *repo.PendingJobsRepo
}

var (
	once    sync.Once
	service *Service
)

func InitialiseService() {
	once.Do(func() {
		service = &Service{
			baseRepo:        repo.GetRepository(),
			pendingJobsRepo: repo.GetPendingJobsRepo(),
		}
	})
}

func GetService() *Service {
	return service
}

func (s *Service) GetPendingJobs(c *gin.Context, request dtos.GetPendingJobsRequest) (interface{}, error) {
	items, err := s.pendingJobsRepo.FetchAllByConditions(c, request)
	if err != nil {
		log.Fatalf("error fetching item: %v", err)
	}

	return items, nil
}
