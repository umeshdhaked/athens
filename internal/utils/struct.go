package utils

import (
	"fmt"
	"reflect"
	"sort"
)

// SortByField sorts an array of structs based on the specified field in ascending or descending order.
// The field must be exported (start with an uppercase letter) for sorting to work.
// The sort order can be "asc" (ascending) or "desc" (descending).
func SortByField(data interface{}, field string, order string) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		panic("SortByField: data must be a slice")
	}

	sort.SliceStable(data, func(i, j int) bool {
		val1 := reflect.Indirect(v.Index(i)).FieldByName(field)
		val2 := reflect.Indirect(v.Index(j)).FieldByName(field)

		switch order {
		case "asc":
			return lessThan(val1, val2)
		case "desc":
			return !lessThan(val1, val2)
		default:
			panic("SortByField: invalid order, must be 'asc' or 'desc'")
		}
	})

	return err
}

// lessThan compares two values and returns true if val1 is less than val2.
// It supports comparison of int, float, and string types.
func lessThan(val1, val2 reflect.Value) bool {
	switch val1.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val1.Int() < val2.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val1.Uint() < val2.Uint()
	case reflect.Float32, reflect.Float64:
		return val1.Float() < val2.Float()
	case reflect.String:
		return val1.String() < val2.String()
	default:
		panic("SortByField: unsupported type for comparison")
	}
}
