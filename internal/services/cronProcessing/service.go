package cronProcessing

import (
	"sync"

	"github.com/umeshdhaked/athens/internal/constants"
	"github.com/umeshdhaked/athens/internal/models"
	"github.com/umeshdhaked/athens/internal/pkg/repo"
	"github.com/umeshdhaked/athens/pkg/dtos"
	"github.com/umeshdhaked/athens/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Service struct {
	baseRepo           repo.IRepository
	cronProcessingRepo repo.ICronProcessingRepo
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
	var items []models.CronProcessing
	err := s.baseRepo.FindMultiplePagination(c,
		&items,
		map[string]interface{}{
			constants.Name:   request.Name,
			constants.Status: request.Status,
		}, dtos.Pagination{
			From: request.From,
			To:   request.To,
		})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	return items, nil
}
