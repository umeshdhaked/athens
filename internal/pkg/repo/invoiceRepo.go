package repo

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

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

func (r InvoiceRepo) CreateInvoice() error {

	return nil
}

func (r InvoiceRepo) GetInvoiceDataFromOrderId() error {

	return nil
}

func (r InvoiceRepo) GetInvoiceDataFromInvoiceId() error {

	return nil
}
