package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rickschubert/ebay-camel-camel-camel/time"
)

type Database struct{}

type trackingInformation struct {
	SearchTerm string       `json:"searchTerm"`
	UserId     string       `json:"userId"`
	Price      float64      `json:"price"`
	MaxTime    time.Minutes `json:"maxTime"`
}

type ItemToTrack struct {
	SearchTerm string `json:"searchTerm"`
	Price      int    `json:"price"`
	UserID     string `json:"userId"`
	MaxTime    int    `json:"maxTime"`
	UUID       string
}

type user struct {
	UUID  string `json:"UUID"`
	Email string `json:"email"`
}

var dynamoClient *dynamodb.DynamoDB

const awsDatabaseRegion = "eu-west-2"

func New() Database {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsDatabaseRegion)},
	))
	dynamoClient = dynamodb.New(sess)
	fmt.Println("Established dynamodb session")
	return Database{}
}

func (Database) GetTracking(trackingId string) trackingInformation {
	result, err := dynamoClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("trackings"),
		Key: map[string]*dynamodb.AttributeValue{
			"UUID": {
				S: aws.String(trackingId),
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Could not retrieve the tracking with UUID %v", trackingId))
	}

	var trackingRetrieved trackingInformation
	err = dynamodbattribute.UnmarshalMap(result.Item, &trackingRetrieved)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal record %v", err))
	}

	return trackingRetrieved
}

func (Database) GetUserEmail(userId string) string {
	result, err := dynamoClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("users"),
		Key: map[string]*dynamodb.AttributeValue{
			"UUID": {
				S: aws.String(userId),
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Could not retrieve the user with UUID %v", userId))
	}

	var userRetrieved user
	errMarsh := dynamodbattribute.UnmarshalMap(result.Item, &userRetrieved)
	if errMarsh != nil {
		panic(fmt.Sprintf("Failed to unmarshal record %v", err))
	}

	return userRetrieved.Email
}

func (Database) CreateTracking(article ItemToTrack) *dynamodb.PutItemOutput {
	av, err := dynamodbattribute.MarshalMap(article)
	if err != nil {
		panic(fmt.Sprintf("There was a problem with unmarshaling the input article %v", article))
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("trackings"),
	}
	itemOutput, errWriting := dynamoClient.PutItem(input)
	if errWriting != nil {
		panic(fmt.Sprintf("An error occured when trying to post the item to DynamoDB: %v", err.Error()))
	}
	return itemOutput
}

func (Database) DeleteTracking(trackingId string) *dynamodb.DeleteItemOutput {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"UUID": {
				S: aws.String(trackingId),
			},
		},
		TableName: aws.String("trackings"),
	}
	deleteItemOutput, err := dynamoClient.DeleteItem(input)
	if err != nil {
		panic(fmt.Sprintf("An error occured when trying to delete the item from DynamoDB: %v", err.Error()))
	}
	return deleteItemOutput
}
