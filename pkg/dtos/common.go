package dtos

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type DbConditions struct {
	Index   string
	PKey    map[string]interface{}
	NonPKey map[string]interface{}
}

type DbUpdateQueryConditions struct {
	Key      map[string]types.AttributeValue
	ToUpdate map[string]types.AttributeValue
}

type DbScanQueryConditions struct {
	Filters map[string]types.AttributeValue
}
