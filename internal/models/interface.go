package models

type IModel interface {
	TableName() string
	GetID() int64
	SetID(int64)
}
