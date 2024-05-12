package contacts

import (
	"log"
	"sync"

	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/internal/utils"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

type Service struct {
	baseRepo     *repo.Repository
	groupRepo    *repo.GroupRepo
	contactsRepo *repo.ContactsRepo
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
	items, err := s.contactsRepo.FetchAllByConditions(c, dtos.GetContactsRequest{
		GroupName: request.GroupName,
	}, "")
	if err != nil {
		log.Fatalf("error fetching item: %v", err)
	}

	err = utils.SortByField(items, "Name", "asc")
	if err != nil {
		log.Fatalf("error sorting items: %v", err)
	}

	return items, nil
}
