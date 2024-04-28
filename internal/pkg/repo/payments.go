package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/fastbiztech/hastinapura/internal/models"
	"log"
	"strconv"
)

var payments *Payments

type Payments struct {
	Repository
}

func newPaymentsRepo(client *dynamodb.Client) {
	payments = &Payments{Repository: Repository{dbClient: client}}
}

func GetPaymentsRepo() *Payments {
	return payments
}

func (p *Payments) GetOrderFromId(ctx context.Context, orderId string) (*models.RzpOrder, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(models.TableRzpOrders),
		KeyConditions: map[string]types.Condition{
			models.ColumnOrderId: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{&types.AttributeValueMemberS{Value: orderId}},
			},
		},
	}

	var resp, err = p.dbClient.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}
	if resp.Count > 0 {
		rzpOrder := []models.RzpOrder{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &rzpOrder); err != nil {
			log.Println(err)
			return nil, err
		}
		return &rzpOrder[0], nil
	}
	return nil, nil

}

func (p *Payments) GetPaymentFromId(ctx context.Context, paymentId string) (*models.RzpPayment, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(models.TableRzpPayments),
		KeyConditions: map[string]types.Condition{
			models.ColumnOrderId: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{&types.AttributeValueMemberS{Value: paymentId}},
			},
		},
	}

	var resp, err = p.dbClient.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}
	if resp.Count > 0 {
		rzpPmt := []models.RzpPayment{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &rzpPmt); err != nil {
			log.Println(err)
			return nil, err
		}
		return &rzpPmt[0], nil
	}
	return nil, nil

}

func (p *Payments) GetRefundFromId(ctx context.Context, refundId string) (*models.RzpRefunds, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(models.TableRzpRefunds),
		KeyConditions: map[string]types.Condition{
			models.ColumnOrderId: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{&types.AttributeValueMemberS{Value: refundId}},
			},
		},
	}

	var resp, err = p.dbClient.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}
	if resp.Count > 0 {
		rzpRfnd := []models.RzpRefunds{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &rzpRfnd); err != nil {
			log.Println(err)
			return nil, err
		}
		return &rzpRfnd[0], nil
	}
	return nil, nil

}

func (p *Payments) GetCreatedOrders(ctx context.Context, createdAtBefore int64) ([]models.RzpOrder, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(models.TableRzpOrders),
		IndexName: aws.String(models.IndexTableRzpOrdersIndexStatus),
		KeyConditions: map[string]types.Condition{
			models.ColumnOrderStatus: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "created"}},
			},
		},
		FilterExpression: aws.String("created_at < :var0"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":var0": &types.AttributeValueMemberN{Value: strconv.FormatInt(createdAtBefore, 10)},
		},
	}

	var resp, err = p.dbClient.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}
	if resp.Count > 0 {
		rzpOrder := []models.RzpOrder{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &rzpOrder); err != nil {
			log.Println(err)
			return nil, err
		}
		return rzpOrder, nil
	}
	return nil, nil

}

func (p *Payments) GetAttemptedOrders(ctx context.Context, createdAtBefore int64) ([]models.RzpOrder, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(models.TableRzpOrders),
		IndexName: aws.String(models.IndexTableRzpOrdersIndexStatus),
		KeyConditions: map[string]types.Condition{
			models.ColumnOrderStatus: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "attempted"}},
			},
		},
		FilterExpression: aws.String("created_at < :var0"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":var0": &types.AttributeValueMemberN{Value: strconv.FormatInt(createdAtBefore, 10)},
		},
	}

	var resp, err = p.dbClient.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}
	if resp.Count > 0 {
		rzpOrder := []models.RzpOrder{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &rzpOrder); err != nil {
			log.Println(err)
			return nil, err
		}
		return rzpOrder, nil
	}
	return nil, nil
}

func (p *Payments) GetOrderList(ctx context.Context, pageLimit int32, lastEvaluatedKey map[string]types.AttributeValueMemberS, status string) ([]models.RzpOrder, map[string]types.AttributeValue, int32, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName:        aws.String(models.TableRzpOrders),
		Limit:            aws.Int32(pageLimit),
		IndexName:        aws.String(models.IndexTableRzpOrdersIndexCreatedAt),
		ScanIndexForward: aws.Bool(false),
		KeyConditions: map[string]types.Condition{
			models.ColumnOrderStatus: {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: status},
				},
			},
			models.ColumnOrderCreatedAt: {
				ComparisonOperator: types.ComparisonOperatorGt,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberN{Value: "0"},
				},
			},
		},
	}

	if lastEvaluatedKey != nil {
		queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
			"id":         &types.AttributeValueMemberS{Value: lastEvaluatedKey["id"].Value},
			"created_at": &types.AttributeValueMemberN{Value: lastEvaluatedKey["created_at"].Value},
			"status":     &types.AttributeValueMemberS{Value: lastEvaluatedKey["status"].Value},
		}
	}

	var resp, err = p.dbClient.Query(ctx, queryInput)
	if err != nil {
		return nil, nil, 0, err
	}
	if resp.Count > 0 {
		rzpOrder := []models.RzpOrder{}
		if err := attributevalue.UnmarshalListOfMaps(resp.Items, &rzpOrder); err != nil {
			log.Println(err)
			return nil, nil, 0, err
		}
		return rzpOrder, resp.LastEvaluatedKey, resp.Count, nil
	}
	return nil, nil, 0, nil
}
