package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var db *dynamodb.DynamoDB

// InitializeDynamoDB create a DynamoDB reference
func InitializeDynamoDB() {
	db = dynamodb.New(session.New(&aws.Config{
		Endpoint: aws.String(os.Getenv("DYNAMODB_ENDPOINT")),
	}))
	log.Println("DynamoDB client initialized.")
}

// GetDynamoDB return the DynamoDB reference
func GetDynamoDB() *dynamodb.DynamoDB {
	return db
}
