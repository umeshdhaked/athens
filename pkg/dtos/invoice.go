package dtos

import "github.com/fastbiztech/hastinapura/internal/models"

type InvoicesRequest struct {
	OrderId string
}

type InvoicesResponse struct {
	*models.Invoice
}
