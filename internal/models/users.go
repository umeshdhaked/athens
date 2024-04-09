package models

const (
	TableUser       = "User"
	TablePromoPhone = "PromoPhonesNo"
	TableOtp        = "Otp"
)

const (
	ColumnMobile             = "Mobile"
	ColumnIsAlreadyContacted = "IsAlreadyContacted"
)

type User struct { //user in dynamoDB
	ID              string `json:"id"`
	Mobile          string `json:"mobile"`
	Hashed_password string `json:"hashed_password"`
	Name            string `json:"name"`
	Role            string `json:"role"`
}

type PromoPhone struct { //promo-phones-no table in dynamoDB
	Mobile             string `json:"mobile"`
	Timestamp          string `json:"timestamp"`
	IsAlreadyContacted string `json:"isAlreadyContacted"`
	Comment            string
}

type Otp struct { // otp table in dynamoDB
	Id     string `json:"id"`
	Mobile string `json:"mobile"`
	Otp    string `json:"otp"`
	Exp    int64  `json:"exp"`
}
