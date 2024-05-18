package models

const (
	TableGroup = "contact_groups"
)

type ContactGroups struct {
	ID          int64  `gorm:"column:id;primary_key" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	UserID      int64  `gorm:"column:user_id" json:"user_id"`
	ColumnNames string `gorm:"column:column_names" json:"column_names"`
	BaseModel
}

func (g *ContactGroups) TableName() string {
	return TableGroup
}

func (g *ContactGroups) GetID() int64 {
	return g.ID
}

func (g *ContactGroups) SetID(id int64) {
	g.ID = id
}
