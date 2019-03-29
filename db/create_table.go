package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var db *dynamodb.DynamoDB

func main() {
	InitializeDB()
	CreateTable()
}

// InitializeDB create a DynamoDB reference
func InitializeDB() {
	fmt.Println("Initializing DynamoDB session ...")

	db = dynamodb.New(session.New(&aws.Config{
		Endpoint: aws.String(os.Getenv("DYNAMODB_ENDPOINT")),
	}))
}

// CreateTable creates the dynamodb table
func CreateTable() {
	fmt.Println("Creating table 'OTF-Workouts' ...")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("otf-workout-id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("otf-workout-id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("OTF-Workouts"),
	}

	_, err := db.CreateTable(input)

	if err != nil {
		fmt.Println("Failed to create table 'OTF-Workouts'. Got error:", err)
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Created the table 'OTF-Workouts' in us-east-1.")
}
