package main

import (
	"context"
	"errors"
	"log"
	"net"
	"strconv"

	pb "../definition"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"google.golang.org/grpc"
)

const port = ":8080"

// WorkoutServiceServer is the server API for WorkoutService service.
type WorkoutServiceServer struct{}

// Item is the model that represent a OTF workout in the database.
type Item struct {
	OTFWorkoutID string `json:"otf-workout-id"`
}

// CreateWorkout creates a workout in DynamoDB
func (s *WorkoutServiceServer) CreateWorkout(ctx context.Context, in *pb.Workout) (*pb.WorkoutResponse, error) {
	log.Printf("Creating workout %s ...\n", in)

	item := &Item{}
	item.OTFWorkoutID = "otf" + strconv.FormatInt(in.GetWorkoutDatetime(), 10)
	// TODO: Store the other workout properties

	av, err := dynamodbattribute.MarshalMap(item)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("OTF-Workouts"),
		ExpressionAttributeNames: map[string]*string{
			"#m": aws.String("otf-workout-id"),
		},
		ConditionExpression: aws.String("attribute_not_exists(#m)"),
	}

	_, err = GetDynamoDB().PutItem(input)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.WorkoutResponse{Message: "Created workout."}, nil
}

// GetWorkout retrieves the workout in DynamoDB by workout id
func (s *WorkoutServiceServer) GetWorkout(ctx context.Context, in *pb.WorkoutRequest) (*pb.Workout, error) {
	log.Printf("Retrieving workout %s ...\n", in)

	result, err := GetDynamoDB().GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("OTF-Workouts"),
		Key: map[string]*dynamodb.AttributeValue{
			"otf-workout-id": {
				S: aws.String(in.GetWorkoutId()),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	item := &Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	if err != nil {
		return nil, err
	}

	if item.OTFWorkoutID == "" {
		return nil, errors.New("OTF workout not found")
	}

	return &pb.Workout{}, nil
}

// UpdateWorkout updates a workout that exist in DynamoDB
func (s *WorkoutServiceServer) UpdateWorkout(ctx context.Context, in *pb.Workout) (*pb.WorkoutResponse, error) {
	log.Printf("Updating workout %s ...\n", in)

	item := &Item{}
	item.OTFWorkoutID = "otf" + strconv.FormatInt(in.GetWorkoutDatetime(), 10)
	// TODO: Update the other workout properties

	av, err := dynamodbattribute.MarshalMap(item)

	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("OTF-Workouts"),
		ExpressionAttributeNames: map[string]*string{
			"#m": aws.String("otf-workout-id"),
		},
		ConditionExpression: aws.String("attribute_exists(#m)"),
	}

	_, err = GetDynamoDB().PutItem(input)

	if err != nil {
		return nil, err
	}

	return &pb.WorkoutResponse{Message: "Updated workout."}, nil
}

// DeleteWorkout deletes a workout from DynamoDB
func (s *WorkoutServiceServer) DeleteWorkout(ctx context.Context, in *pb.WorkoutRequest) (*pb.WorkoutResponse, error) {
	log.Printf("Deleting workout %s ...\n", in)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"otf-workout-id": {
				S: aws.String(in.GetWorkoutId()),
			},
		},
		TableName: aws.String("OTF-Workouts"),
		ExpressionAttributeNames: map[string]*string{
			"#m": aws.String("otf-workout-id"),
		},
		ConditionExpression: aws.String("attribute_exists(#m)"),
	}

	_, err := GetDynamoDB().DeleteItem(input)

	if err != nil {
		return nil, err
	}

	return &pb.WorkoutResponse{Message: "Deleted workout."}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to startup grpc server: %s\n", err)
	}
	s := grpc.NewServer()
	pb.RegisterWorkoutServiceServer(s, &WorkoutServiceServer{})
	log.Println("Initializing DynamoDB client ...")
	InitializeDynamoDB()
	log.Printf("Starting grpc server on port %s ...", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to startup grpc server: %s\n", err)
	}
}
