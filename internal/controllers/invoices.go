package controllers

import (
	"net/http"

	"github.com/umeshdhaked/athens/internal/services/invoices"
	"github.com/umeshdhaked/athens/pkg/dtos"
	"github.com/gin-gonic/gin"
)

func HandleInvoice(ctx *gin.Context) {
	var invoiceReq dtos.InvoicesRequest
	if err := ctx.ShouldBindJSON(&invoiceReq); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var inv *invoices.InvoiceService = invoices.GeInvoiceService()
	invResp, err := inv.GetInvoiceFromOrderId(ctx, &invoiceReq)
	if nil != err {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, invResp)
}
