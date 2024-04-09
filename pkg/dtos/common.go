package dtos

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type DbQueryInputConditions struct {
	Index   string
	PKey    map[string]interface{}
	NonPKey map[string]interface{}
}

type DbFilterQueryConditions struct {
	Filters map[string]types.AttributeValue
}

type DbUpdateQueryConditions struct {
	Key      map[string]types.AttributeValue
	ToUpdate map[string]types.AttributeValue
}
