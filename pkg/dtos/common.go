package dtos

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type DbQueryInputConditions struct {
	Index   string
	PKey    map[string]interface{}
	NonPKey map[string]interface{}
	Limit   int
}

type DbScanQueryConditions struct {
	ExclusiveStartKey map[string]types.AttributeValue
	Filters           map[string]types.AttributeValue
	Limit             int
}

type DbUpdateQueryConditions struct {
	Key      map[string]types.AttributeValue
	ToUpdate map[string]types.AttributeValue
}
