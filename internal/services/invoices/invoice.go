package invoices

import (
	"fmt"
	"sync"

	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/fastbiztech/hastinapura/internal/pkg/repo"
	"github.com/fastbiztech/hastinapura/pkg/dtos"
	"github.com/gin-gonic/gin"
)

var (
	once           sync.Once
	invoiceService *InvoiceService
)

type InvoiceService struct {
	baseRepo    repo.IRepository
	invoiceRepo repo.IInvoiceRepo
}

func GeInvoiceService() *InvoiceService {
	return invoiceService
}

func NewInvoiceService(invoiceRepo repo.IInvoiceRepo) {
	once.Do(func() {
		invoiceService = &InvoiceService{
			baseRepo:    repo.GetRepository(),
			invoiceRepo: invoiceRepo,
		}
	})
}

func (i *InvoiceService) GetInvoiceFromOrderId(ctx *gin.Context, invoiceReq *dtos.InvoicesRequest) (*dtos.InvoicesResponse, error) {
	invoice := &models.Invoice{}
	err := i.baseRepo.Find(ctx, invoice, map[string]interface{}{
		models.SQLColumnInvoiceOrderId: invoiceReq.OrderId,
	})
	if nil != err {
		return nil, err
	}
	rzpOrder := &models.Payments{}
	err = i.baseRepo.Find(ctx, rzpOrder, map[string]interface{}{
		models.SQLColumnInvoiceOrderId: invoiceReq.OrderId,
	})
	return &dtos.InvoicesResponse{
		Invoice:     invoice,
		InvoiceId:   fmt.Sprintf("INVFBT%06d", invoice.ID),
		RzpOrderDBO: rzpOrder,
	}, nil
}
