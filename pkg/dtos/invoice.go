package dtos

import "github.com/umeshdhaked/athens/internal/models"

type InvoicesRequest struct {
	OrderId string
}

type InvoicesResponse struct {
	InvoiceId   string
	Invoice     *models.Invoice
	RzpOrderDBO *models.Payments
}
