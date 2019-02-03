package main

import (
	"context"
	"log"
	"net"

	pb "../definition"
	"google.golang.org/grpc"
)

const port = ":8080"

// WorkoutServiceServer is the server API for WorkoutService service.
type WorkoutServiceServer struct{}

// CreateWorkout creates a workout log in DynamoDB
func (s *WorkoutServiceServer) CreateWorkout(ctx context.Context, in *pb.Workout) (*pb.WorkoutResponse, error) {
	log.Printf("Creating workout %s ...\n", in)
	return &pb.WorkoutResponse{Message: "Created workout."}, nil
}

// GetWorkout retrieves the workout log in DynamoDB by workout id
func (s *WorkoutServiceServer) GetWorkout(ctx context.Context, in *pb.WorkoutRequest) (*pb.Workout, error) {
	log.Printf("Retrieving workout %s ...\n", in)
	return &pb.Workout{}, nil
}

// UpdateWorkout updates a workout log that exist in DynamoDB
func (s *WorkoutServiceServer) UpdateWorkout(ctx context.Context, in *pb.Workout) (*pb.WorkoutResponse, error) {
	log.Printf("Updating workout %s ...\n", in)
	return &pb.WorkoutResponse{Message: "Updated workout."}, nil
}

// DeleteWorkout deletes a workout log from DynamoDB
func (s *WorkoutServiceServer) DeleteWorkout(ctx context.Context, in *pb.WorkoutRequest) (*pb.WorkoutResponse, error) {
	log.Printf("Deleting workout %s ...\n", in)
	return &pb.WorkoutResponse{Message: "Deleted workout."}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to startup grpc server: %s\n", err)
	}
	s := grpc.NewServer()
	pb.RegisterWorkoutServiceServer(s, &WorkoutServiceServer{})
	log.Printf("Starting grpc server on port %s ...", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to startup grpc server: %s\n", err)
	}
}
