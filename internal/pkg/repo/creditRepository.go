package repo

//type CreditsRepo struct {
//	client *dynamodb.Client
//}
//
//func NewCreditsRepo(client *dynamodb.Client) *CreditsRepo {
//	return &CreditsRepo{client: client}
//}

//func (c *CreditsRepo) CreateUserCredit(ctx *gin.Context, credit *models.Credits) error {
//	item, _ := attributevalue.MarshalMap(credit)
//	params := &dynamodb.PutItemInput{
//		TableName: aws.String("credit"),
//		Item:      item,
//	}
//
//	output, er := c.client.PutItem(ctx, params)
//	log.Println(output)
//	return er
//}
