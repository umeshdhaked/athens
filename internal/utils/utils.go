package utils

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/umeshdhaked/athens/internal/constants"
)

// IsEmpty will check for given data is empty as per the go documentation
func IsEmpty(val interface{}) bool {
	if val == nil {
		return true
	}

	reflectVal := reflect.ValueOf(val)

	switch reflectVal.Kind() {
	case reflect.Int:
		return val.(int) == 0

	case reflect.Int64:
		return val.(int64) == 0

	case reflect.Float32, reflect.Float64:
		return reflectVal.Float() == 0

	case reflect.String:
		return strings.TrimSpace(val.(string)) == ""

	case reflect.Map:
		fallthrough
	case reflect.Slice:
		return reflectVal.IsNil() || reflectVal.Len() == 0

	case reflect.Interface, reflect.Ptr:
		if reflectVal.IsNil() {
			return true
		}
		return IsEmpty(reflectVal.Elem().Interface())

	case reflect.Struct:
		copyStruct := reflect.New(reflect.TypeOf(val)).Elem().Interface()
		if reflect.DeepEqual(val, copyStruct) {
			return true
		}
	}

	return false
}

func GetFilePath(name string) string {
	dirPath, err := os.Getwd()
	if err != nil {
		dirPath = fmt.Sprintf(constants.BasePath, os.Getenv("GOPATH"))
	}
	return fmt.Sprintf("%s/%s", dirPath, name)
}

func Ternary(condition bool, trueValue, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}

func IsEmptyAttributeValue(value types.AttributeValue) bool {
	var valueMap interface{}
	err := attributevalue.Unmarshal(value, &valueMap)
	if err != nil {
		return false
	}

	if IsEmpty(valueMap) {
		return true
	}

	return false
}

// Function to check if a string is present in a slice of strings
func Contains(arr []string, target string) bool {
	// Use the Contains method from the strings package
	for _, str := range arr {
		if str == target {
			return true
		}
	}
	return false
}
