package contacts

import (
	"sync"

	"github.com/umeshdhaked/athens/internal/constants"
	"github.com/umeshdhaked/athens/internal/models"
	"github.com/umeshdhaked/athens/internal/pkg/repo"
	"github.com/umeshdhaked/athens/internal/utils"
	"github.com/umeshdhaked/athens/pkg/dtos"
	"github.com/umeshdhaked/athens/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Service struct {
	baseRepo     repo.IRepository
	groupRepo    repo.IGroupRepo
	contactsRepo repo.IContactsRepo
}

var (
	once    sync.Once
	service *Service
)

func InitialiseService() {
	once.Do(func() {
		service = &Service{
			baseRepo:     repo.GetRepository(),
			groupRepo:    repo.GetGroupRepo(),
			contactsRepo: repo.GetContactsRepo(),
		}
	})
}

func GetService() *Service {
	return service
}

func (s *Service) GetContacts(c *gin.Context, request dtos.GetGroupContactsRequest) (interface{}, error) {
	var items []models.Contacts
	err := s.baseRepo.FindMultiplePagination(c,
		&items,
		map[string]interface{}{
			constants.GroupName: request.GroupName,
		},
		dtos.Pagination{})
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	// todo: add order by.
	err = utils.SortByField(items, "Name", "asc")
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("error sorting items:")
	}

	return items, nil
}
