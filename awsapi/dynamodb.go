package awsapi

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

func NewDynamoDbAPI() *DynamoDbAPI {

	client := dynamodb.New(LoadAWSConfig())

	return &DynamoDbAPI{
		client: client,
	}
}

type DynamoDbAPI struct {
	client *dynamodb.Client
}

func (api *DynamoDbAPI) GetItem(tableName, keyName, keyVal string, item interface{}) error {

	key := map[string]dynamodb.AttributeValue{
		keyName: {
			S: aws.String(keyVal),
		},
	}

	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(tableName),
	}

	req := api.client.GetItemRequest(input)

	resp, err := req.Send(context.Background())
	if err != nil {
		return err
	}
	return dynamodbattribute.UnmarshalMap(resp.Item, item)
}

func (api *DynamoDbAPI) PutItem(tableName string, itemVal interface{}) error {

	item, err := dynamodbattribute.MarshalMap(itemVal)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}

	req := api.client.PutItemRequest(input)

	_, err = req.Send(context.Background())
	return err
}
