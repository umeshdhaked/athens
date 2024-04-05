package dbo

type User struct { //user_table in dynamoDB
	Id              string `json:"id"`
	Mobile          string `json:"mobile"`
	Hashed_password string `json:"hashed_password"`
	Name            string `json:"name"`
	Role            string `json:"role"`
}

type PromoPhone struct { //promo_phones_no table in dynamoDB
	Mobile             string `json:"mobile"`
	Timestamp          string `json:"timestamp"`
	IsAlreadyContacted string `json:"is_already_contacted"`
	Comment            string
}

type Otp struct { // otp table in dynamoDB
	Id     string `json:"id"`
	Mobile string `json:"mobile"`
	Otp    string `json:"otp"`
	Exp    int64  `json:"exp"`
}
