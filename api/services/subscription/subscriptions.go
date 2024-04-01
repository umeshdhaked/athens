package subscription

import (
	"errors"
	"log"
	"time"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/fastbiztech/hastinapura/pkg/models/dbo"
	"github.com/fastbiztech/hastinapura/pkg/models/requests"
	"github.com/fastbiztech/hastinapura/pkg/models/responses"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionService struct {
	svc *dynamodb.DynamoDB
}

func NewSubscriptionService(svc *dynamodb.DynamoDB) *SubscriptionService {
	return &SubscriptionService{svc: svc}
}

func (s *SubscriptionService) CreateNewPricingSystem(ctx *gin.Context, pricing *requests.PricingRequest) (*responses.PricingResponse, error) {

	if role, exists := ctx.Params.Get("role"); !exists {
		return nil, errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return nil, errors.New("only admin user allowed to add pricing")
	}

	if "" == pricing.Category || "" == pricing.SubCatgory || "" == pricing.Type || 0 == pricing.Rates {
		return nil, errors.New("invalid input")
	}
	if !(pricing.Category == "SMS" || pricing.Category == "EMAIL" || pricing.Category == "WHATSAPP") {
		return nil, errors.New("inavaid category")
	}
	if pricing.Category == "SMS" && pricing.SubCatgory != "PROMOTIONAL" && pricing.SubCatgory != "TRANSACTIONAL" {
		return nil, errors.New("inavaid sub_category for SMS")
	}

	// search in DB if default exists.
	if pricing.Type == "DEFAULT" {
		var queryInput = &dynamodb.QueryInput{
			TableName:              aws.String("pricing"),
			IndexName:              aws.String("category-index"),
			KeyConditionExpression: aws.String("category = :var0"),
			FilterExpression:       aws.String("sub_category= :var1 and pricing_type = :var2 and pricing_state = :var3"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":var0": {S: aws.String(pricing.Category)},
				":var1": {S: aws.String(pricing.SubCatgory)},
				":var2": {S: aws.String("DEFAULT")},
				":var3": {S: aws.String("ACTIVE")},
			}}
		var resp1, err1 = s.svc.Query(queryInput)
		if err1 != nil {
			return nil, err1
		}
		if *resp1.Count > 0 {
			return nil, errors.New("default pricing already exists for given category, subcategory")
		}
	}

	// save pricing to db.
	obj := dbo.Pricing{
		Id:           uuid.New().String(),
		Category:     pricing.Category,
		SubCatgory:   pricing.SubCatgory,
		PricingType:  pricing.Type,
		Rates:        pricing.Rates,
		PricingState: "ACTIVE",
		CreatedAt:    time.Now().Unix(),
	}

	item, _ := dynamodbattribute.MarshalMap(obj)
	params := &dynamodb.PutItemInput{
		TableName: aws.String("pricing"),
		Item:      item,
	}

	req, output := s.svc.PutItemRequest(params)
	fmt.Print(output)
	er := req.Send()
	if er != nil {
		return nil, errors.Join(er, errors.New("FAILED TO MAKE API CALL TO DYNAMO for pricing"))
	}

	return &responses.PricingResponse{Id: obj.Id}, nil
}

func (s *SubscriptionService) FetchAllActivePricingModel(ctx *gin.Context) ([]*responses.PricingResponse, error) {
	if role, exists := ctx.Params.Get("role"); !exists {
		return nil, errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return nil, errors.New("only admin user allowed to add pricing")
	}
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("pricing"),
		IndexName: aws.String("pricing_state-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"pricing_state": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String("ACTIVE"),
					},
				},
			},
		},
	}

	var resp1, err1 = s.svc.Query(queryInput)
	if err1 != nil {
		fmt.Println(err1)
		return nil, err1
	}
	pricing := []dbo.Pricing{}
	if err := dynamodbattribute.UnmarshalListOfMaps(resp1.Items, &pricing); err != nil {
		fmt.Println(err)
	}

	resp := []*responses.PricingResponse{}
	for _, p := range pricing {
		resp = append(resp, &responses.PricingResponse{Id: p.Id, Category: p.Category, SubCatgory: p.SubCatgory, Type: p.PricingType, Rates: p.Rates, Status: p.PricingState})
	}

	return resp, nil
}

func (s *SubscriptionService) AddDefaultSubscriptionToUser(ctx *gin.Context, subReq *requests.UserDefaultSubscriptionRequest) error {
	if role, exists := ctx.Params.Get("role"); !exists {
		return errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return errors.New("only admin user allowed to add pricing")
	}

	// get User
	mobile := subReq.UserMobile
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("user_table"),
		IndexName: aws.String("mobile-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"mobile": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(mobile),
					},
				},
			},
		},
	}
	var userResp, err = s.svc.Query(queryInput)
	users := []dbo.User{}
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		if err := dynamodbattribute.UnmarshalListOfMaps(userResp.Items, &users); err != nil {
			fmt.Println(err)
		}
		log.Println(users)

	}

	// get default Pricing models
	var queryInput1 = &dynamodb.QueryInput{
		TableName:              aws.String("pricing"),
		IndexName:              aws.String("pricing_state-index"),
		KeyConditionExpression: aws.String("pricing_state = :var0"),
		FilterExpression:       aws.String("pricing_type= :var1"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":var0": {S: aws.String("ACTIVE")},
			":var1": {S: aws.String("DEFAULT")},
		},
	}

	var pricingResp, err1 = s.svc.Query(queryInput1)
	if err1 != nil {
		fmt.Println(err1)
		return err1
	}
	defaultPricings := []dbo.Pricing{}
	if err := dynamodbattribute.UnmarshalListOfMaps(pricingResp.Items, &defaultPricings); err != nil {
		fmt.Println(err)
	}

	admin, _ := ctx.Params.Get("id")
	// create Subscriptions to USER
	writeRequests := []*dynamodb.WriteRequest{}
	for _, p := range defaultPricings {
		userSubsDto := dbo.UserSubscription{
			Id:        uuid.New().String(),
			PricingId: p.Id,
			UserId:    users[0].Id,
			Type:      p.Category,
			SubType:   p.SubCatgory,
			Status:    "ACTIVE",
			AddedBy:   admin,
			CreatedAt: time.Now().Unix(),
		}
		item, _ := dynamodbattribute.MarshalMap(userSubsDto)
		writeRequests = append(writeRequests,
			&dynamodb.WriteRequest{PutRequest: &dynamodb.PutRequest{Item: item}})
	}

	batchWrite := dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"user_subscriptions": writeRequests,
		},
	}
	output, er := s.svc.BatchWriteItemWithContext(ctx, &batchWrite)
	fmt.Println(output)
	if er != nil {
		return errors.Join(er, errors.New("FAILED TO MAKE API CALL TO DYNAMO for user_subscriptions"))
	}

	return nil
}

func (s *SubscriptionService) AddSubscriptionToUser(ctx *gin.Context, subReq *requests.UserSubscriptionRequest) error {
	if role, exists := ctx.Params.Get("role"); !exists {
		return errors.New("internal server error, user role not found")
	} else if "admin" != role {
		return errors.New("only admin user allowed to add pricing")
	}

	// get User
	mobile := subReq.UserMobile
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("user_table"),
		IndexName: aws.String("mobile-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"mobile": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(mobile),
					},
				},
			},
		},
	}
	var userResp, err = s.svc.Query(queryInput)
	users := []dbo.User{}
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		if err := dynamodbattribute.UnmarshalListOfMaps(userResp.Items, &users); err != nil {
			fmt.Println(err)
		}
		log.Println(users)

	}

	// get default Pricing models
	var queryInput1 = &dynamodb.QueryInput{
		TableName:              aws.String("pricing"),
		IndexName:              aws.String("pricing_state-index"),
		KeyConditionExpression: aws.String("pricing_state= :var1"),
		FilterExpression:       aws.String("id = :var0"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":var0": {S: aws.String(subReq.PricingId)},
			":var1": {S: aws.String("ACTIVE")},
		},
	}

	var pricingResp, err1 = s.svc.Query(queryInput1)
	if err1 != nil {
		fmt.Println(err1)
		return err1
	}
	pricings := []dbo.Pricing{}
	if err := dynamodbattribute.UnmarshalListOfMaps(pricingResp.Items, &pricings); err != nil {
		fmt.Println(err)
	}

	// create Subscriptions to USER
	admin, _ := ctx.Params.Get("id")

	userSubsDto := dbo.UserSubscription{
		Id:        uuid.New().String(),
		PricingId: pricings[0].Id,
		UserId:    users[0].Id,
		Type:      pricings[0].Category,
		SubType:   pricings[0].SubCatgory,
		Status:    "ACTIVE",
		AddedBy:   admin,
		CreatedAt: time.Now().Unix(),
	}

	item, _ := dynamodbattribute.MarshalMap(userSubsDto)
	params := &dynamodb.PutItemInput{
		TableName: aws.String("user_subscriptions"),
		Item:      item,
	}

	req, output := s.svc.PutItemRequest(params)
	fmt.Print(output)
	er := req.Send()
	if er != nil {
		return errors.Join(er, errors.New("FAILED TO MAKE API CALL TO DYNAMO for user_subscriptions"))
	}
	return nil
}

func (s *SubscriptionService) FetchAllActiveSubscriptionsForUser(ctx *gin.Context, req *requests.FetchSubscriptionRequest) ([]*responses.SubscriptionResponse, error) {
	if role, exists := ctx.Params.Get("role"); !exists {
		return nil, errors.New("internal server error, user role not found")
	} else if role != "admin" {
		return nil, errors.New("only admin user allowed to add pricing")
	}

	// get User
	mobile := req.UserMobile
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("user_table"),
		IndexName: aws.String("mobile-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"mobile": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(mobile),
					},
				},
			},
		},
	}
	var userResp, err = s.svc.Query(queryInput)
	users := []dbo.User{}
	if err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		if err := dynamodbattribute.UnmarshalListOfMaps(userResp.Items, &users); err != nil {
			fmt.Println(err)
		}
		log.Println(users)

	}

	queryInput = &dynamodb.QueryInput{
		TableName: aws.String("user_subscriptions"),
		IndexName: aws.String("user_id-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"user_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(users[0].Id),
					},
				},
			},
		},
	}

	var resp1, err1 = s.svc.Query(queryInput)
	if err1 != nil {
		fmt.Println(err1)
		return nil, err1
	}
	subscriptions := []dbo.UserSubscription{}
	if err := dynamodbattribute.UnmarshalListOfMaps(resp1.Items, &subscriptions); err != nil {
		fmt.Println(err)
	}

	resp := []*responses.SubscriptionResponse{}
	for _, s := range subscriptions {
		resp = append(resp, &responses.SubscriptionResponse{
			Id:        s.Id,
			PricingId: s.PricingId,
			UserId:    s.UserId,
			Type:      s.Type,
			SubType:   s.SubType,
			Status:    s.Status,
			AddedBy:   s.AddedBy,
			CreatedAt: s.CreatedAt,
			DeletedAt: s.DeletedAt,
		})
	}

	return resp, nil
}

// TODO: add/update credit api

func chargeUser(userId string, typ string, subtype string) {
	// get subscription (recent one for typ,subtype,userId)
	// get pricing from subscription
	// fetch credit for user_id
	// update credit from user_id
}
