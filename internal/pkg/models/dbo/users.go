package dbo

type User struct { //user_table in dynamoDB
	Id              string `json:"id" dynamodbav:"id"`
	Mobile          string `json:"mobile" dynamodbav:"mobile"`
	Hashed_password string `json:"hashed_password" dynamodbav:"hashed_password"`
	Name            string `json:"name" dynamodbav:"name"`
	Role            string `json:"role" dynamodbav:"role"`
}

type PromoPhone struct { //promo_phones_no table in dynamoDB
	Mobile             string `json:"mobile" dynamodbav:"mobile"`
	Timestamp          string `json:"timestamp" dynamodbav:"timestamp"`
	IsAlreadyContacted string `json:"is_already_contacted" dynamodbav:"is_already_contacted"`
	Comment            string
}

type Otp struct { // otp table in dynamoDB
	Id     string `json:"id" dynamodbav:"id"`
	Mobile string `json:"mobile" dynamodbav:"mobile"`
	Otp    string `json:"otp" dynamodbav:"otp"`
	Exp    int64  `json:"exp" dynamodbav:"exp"`
}
