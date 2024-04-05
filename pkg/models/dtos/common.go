package dtos

type DbConditions struct {
	Index   string
	PKey    map[string]interface{}
	NonPKey map[string]interface{}
}
