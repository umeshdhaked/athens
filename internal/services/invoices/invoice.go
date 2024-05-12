package invoices

import (
	"sync"

	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

var (
	once           sync.Once
	invoiceService *InvoiceService
)

type InvoiceService struct {
	invoiceRepo *repo.InvoiceRepo
}

func GeInvoiceService() *InvoiceService {
	return invoiceService
}

func NewInvoiceService(invoiceRepo *repo.InvoiceRepo) {
	once.Do(func() {
		invoiceService = &InvoiceService{invoiceRepo: invoiceRepo}
	})
}

func (i *InvoiceService) GetInvoiceFromOrderId(ctx *gin.Context, invoiceReq *dtos.InvoicesRequest) (*dtos.InvoicesResponse, error) {
	invoice, err := i.invoiceRepo.GetInvoiceDataFromOrderId(ctx, invoiceReq.OrderId)
	if nil != err {
		return nil, err
	}
	return &dtos.InvoicesResponse{Invoice: invoice}, nil
}
