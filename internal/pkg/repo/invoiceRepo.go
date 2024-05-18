package repo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IInvoiceRepo interface {
}

var invoiceRepo IInvoiceRepo

type MysqlInvoiceRepo struct {
	IRepository
}

func newInvoiceRepo(mysqlDB *gorm.DB, dynamoDB *dynamodb.Client) {
	invoiceRepo = &MysqlInvoiceRepo{
		IRepository: &MysqlRepository{
			db: mysqlDB,
		},
	}
}

func GetInvoiceRepo() IInvoiceRepo {
	return invoiceRepo
}

func (r MysqlInvoiceRepo) GetEmptyInvoice(ctx *gin.Context) (*models.Invoice, error) {
	// todo implement
	//invoicecounter, err := r.GetInvoiceFromInvoiceId(ctx, "invoicecounter")
	//if nil != err {
	//	return nil, err
	//}
	//if nil == invoicecounter {
	//	invoicecounter = &models.Invoice{
	//		ID:            uuid.New().String(),
	//		InvoiceId:     "invoicecounter",
	//		InvoiceNumber: 0,
	//		Status:        "INVALID",
	//	}
	//	invMap, err := attributevalue.MarshalMap(invoicecounter)
	//	if err != nil {
	//		return nil, err
	//	}
	//	err = r.Repository.CreateItem(ctx, models.TableInvoices, invMap)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	//newInvoiceNumber := invoicecounter.InvoiceNumber
	//// save +1 invoice counter.
	//invoicecounter.InvoiceNumber = newInvoiceNumber + 1
	//invoiceCounterMp, err := attributevalue.MarshalMap(invoicecounter)
	//if err != nil {
	//	return nil, err
	//}
	//er := r.Repository.CreateItem(ctx, models.TableInvoices, invoiceCounterMp)
	//if er != nil {
	//	return nil, err
	//}
	////create new invoice in db and return
	//invoiceNum := fmt.Sprintf("INVFBT%06d", newInvoiceNumber)
	//
	//newInvoice := &models.Invoice{
	//	ID:            uuid.New().String(),
	//	InvoiceNumber: newInvoiceNumber,
	//	InvoiceId:     invoiceNum,
	//	Status:        "CANCELLED",
	//}
	//invoiceMp, err := attributevalue.MarshalMap(newInvoice)
	//if err != nil {
	//	return nil, err
	//}
	//er = r.Repository.CreateItem(ctx, models.TableInvoices, invoiceMp)
	//if er != nil {
	//	return nil, er
	//}
	//return newInvoice, nil
	return nil, nil
}
