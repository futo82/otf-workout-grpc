package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	pb "../definition"
	"google.golang.org/grpc"
)

const address = "localhost:8080"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to the grpc server: %s", err)
	}
	defer conn.Close()
	c := pb.NewWorkoutServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println("Parsing file 'workout_log.csv' ...")

	csvFile, _ := os.Open("workout_log.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, err := reader.Read()
		if err != nil {
			break
		}
		log.Printf("Creating workout %s ...\n", line)
		r, err := c.CreateWorkout(ctx, &pb.Workout{
			WorkoutDatetime: func() int64 {
				layout := "2006-01-02"
				time, err := time.Parse(layout, line[0])
				if err != nil {
					return 0
				}
				return time.UnixNano()
			}(),
			Coach: line[1],
			GrayZone: func() int32 {
				value, _ := strconv.ParseInt(line[2], 10, 32)
				return int32(value)
			}(),
			BlueZone: func() int32 {
				value, _ := strconv.ParseInt(line[3], 10, 32)
				return int32(value)
			}(),
			GreenZone: func() int32 {
				value, _ := strconv.ParseInt(line[4], 10, 32)
				return int32(value)
			}(),
			OrangeZone: func() int32 {
				value, _ := strconv.ParseInt(line[5], 10, 32)
				return int32(value)
			}(),
			RedZone: func() int32 {
				value, _ := strconv.ParseInt(line[6], 10, 32)
				return int32(value)
			}(),
			CaloriesBurned: func() int32 {
				value, _ := strconv.ParseInt(line[7], 10, 32)
				return int32(value)
			}(),
			SplatPoints: func() int32 {
				value, _ := strconv.ParseInt(line[8], 10, 32)
				return int32(value)
			}(),
			AverageHeartRate: func() int32 {
				value, _ := strconv.ParseInt(line[9], 10, 32)
				return int32(value)
			}(),
		})

		if err != nil {
			log.Fatalf("Response: %s\n", err)
		}
		log.Printf("Response: %s\n", r.Message)
	}
}
