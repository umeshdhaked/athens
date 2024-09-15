package models

const (
	TableKyc = "kyc"
)

type Kyc struct {
	ID           int64
	UserId       string
	UserName     string
	Mobile       string
	DocumentType string
	KycDocLink   string
	PhotoLink    string
	IsVerified   string
	Comment      string
	BaseModel
}

func (o *Kyc) TableName() string {
	return TableKyc
}

func (o *Kyc) GetID() int64 {
	return o.ID
}

func (o *Kyc) SetID(id int64) {
	o.ID = id
}
