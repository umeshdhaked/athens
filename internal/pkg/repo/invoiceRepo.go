package repo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var invoiceRepo *InvoiceRepo

type InvoiceRepo struct {
	Repository
}

func newInvoiceRepo(client *dynamodb.Client) {
	invoiceRepo = &InvoiceRepo{Repository: Repository{dbClient: client}}
}

func GetInvoiceRepo() *InvoiceRepo {
	return invoiceRepo
}

func (r InvoiceRepo) GetEmptyInvoice(ctx *gin.Context) (*models.Invoice, error) {
	invoicecounter, err := r.GetInvoiceFromInvoiceId(ctx, "invoicecounter")
	if nil != err {
		return nil, err
	}
	if nil == invoicecounter {
		invoicecounter = &models.Invoice{
			ID:            uuid.New().String(),
			InvoiceId:     "invoicecounter",
			InvoiceNumber: 0,
			Status:        "INVALID",
			OrderId:       "dummyOrder",
		}
		invMap, err := attributevalue.MarshalMap(invoicecounter)
		if err != nil {
			return nil, err
		}
		err = r.Repository.CreateItem(ctx, models.TableInvoices, invMap)
		if err != nil {
			return nil, err
		}
	}
	newInvoiceNumber := invoicecounter.InvoiceNumber
	// save +1 invoice counter.
	invoicecounter.InvoiceNumber = newInvoiceNumber + 1
	invoiceCounterMp, err := attributevalue.MarshalMap(invoicecounter)
	if err != nil {
		return nil, err
	}
	er := r.Repository.CreateItem(ctx, models.TableInvoices, invoiceCounterMp)
	if er != nil {
		return nil, err
	}
	//create new invoice in db and return
	invoiceNum := fmt.Sprintf("INVFBT%06d", newInvoiceNumber)

	newInvoice := &models.Invoice{
		ID:            uuid.New().String(),
		InvoiceNumber: newInvoiceNumber,
		InvoiceId:     invoiceNum,
		Status:        "CANCELLED",
		OrderId:       "dummyOrder",
	}
	invoiceMp, err := attributevalue.MarshalMap(newInvoice)
	if err != nil {
		return nil, err
	}
	er = r.Repository.CreateItem(ctx, models.TableInvoices, invoiceMp)
	if er != nil {
		return nil, err
	}
	return newInvoice, err
}

func (r *InvoiceRepo) GetInvoiceFromInvoiceId(ctx context.Context, invoiceId string) (*models.Invoice, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(models.TableInvoices),
		IndexName: aws.String(models.IndexTableInvoicesIndexInvoiceId),
		KeyConditions: map[string]types.Condition{
			models.ColumnInvoiceId: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: invoiceId},
				},
			},
		},
	}
	var resp, err = r.dbClient.Query(ctx, queryInput)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if resp.Count > 0 {
		invoices := []models.Invoice{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &invoices); err != nil {
			fmt.Println(err)
		}
		return &invoices[0], nil
	} else {
		return nil, nil
	}
}

func (r InvoiceRepo) GetInvoiceDataFromOrderId(ctx *gin.Context, orderId string) (*models.Invoice, error) {

	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(models.TableInvoices),
		IndexName: aws.String(models.IndexTableInvoicesIndexOrderId),
		KeyConditions: map[string]types.Condition{
			models.ColumnOdrId: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: orderId},
				},
			},
		},
	}
	var resp, err = r.dbClient.Query(ctx, queryInput)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if resp.Count > 0 {
		invoices := []models.Invoice{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &invoices); err != nil {
			fmt.Println(err)
		}
		return &invoices[0], nil
	} else {
		return nil, nil
	}
}
