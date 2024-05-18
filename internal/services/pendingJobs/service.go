package pendingJobs

import (
	"github.com/fastbiztech/hastinapura/pkg/logger"
	"sync"

	"github.com/fastbiztech/hastinapura/internal/constants"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

type Service struct {
	baseRepo        repo.IRepository
	pendingJobsRepo repo.IPendingJobsRepo
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
	var items []models.PendingJobs
	err := s.baseRepo.FindMultiplePagination(c, &items, map[string]interface{}{
		constants.Name:   request.Name,
		constants.Status: request.Status,
		constants.Type:   request.Type,
	}, dtos.Pagination{
		From: request.From,
		To:   request.To,
	})
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("error fetching item:")
	}

	return items, nil
}
