package models

const (
	TableUser       = "user"
	TablePromoPhone = "promo_phones"
	TableOtp        = "otp"
)

const (
	ColumnOtpMobile          = "Mobile"
	ColumnIsAlreadyContacted = "is_already_contacted"

	ColumnUserMobile = "mobile"

	ColumnPromoPhoneMobile = "mobile"
)

type User struct {
	ID              int64  `json:"id"`
	Mobile          string `json:"mobile"`
	Hashed_password string `json:"hashed_password"`
	Name            string `json:"name"`
	Role            string `json:"role"`
	KycDone         string `json:"kyc_done"`
	BaseModel
}

type PromoPhone struct {
	ID                 int64  `json:"id"`
	Mobile             string `json:"mobile"`
	IsAlreadyContacted string `json:"is_already_contacted"`
	Comment            string `json:"comment"`
	BaseModel
}

type Otp struct {
	ID     int64  `json:"id"`
	Mobile string `json:"mobile"`
	Otp    string `json:"otp"`
	Exp    int64  `json:"exp"`
	BaseModel
}

func (o *Otp) TableName() string {
	return TableOtp
}

func (o *Otp) GetID() int64 {
	return o.ID
}

func (o *Otp) SetID(id int64) {
	o.ID = id
}

func (o *User) TableName() string {
	return TableUser
}

func (o *User) GetID() int64 {
	return o.ID
}

func (o *User) SetID(id int64) {
	o.ID = id
}

func (o *PromoPhone) TableName() string {
	return TablePromoPhone
}

func (o *PromoPhone) GetID() int64 {
	return o.ID
}

func (o *PromoPhone) SetID(id int64) {
	o.ID = id
}
