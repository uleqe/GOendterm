package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"com.grpc.nurs/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewCalculatorServiceClient(conn)
	Task1(c)
	Task2(c)
}


func Task2(c greetpb.CalculatorServiceClient) {

	requests := []*greetpb.AvgRequest{
		{
			Number: 1,
		},
		{
			Number: 2,
		},
		{
			Number: 3,
		},
		{
			Number: 4,
		},
	}

	ctx := context.Background()
	stream, err := c.Avg(ctx)
	if err != nil {
		log.Fatalf("error while calling Average: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending number: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from Average: %v", err)
	}
	fmt.Printf("Average Response: %v\n", res)
}




func Task1(c greetpb.CalculatorServiceClient) {
	ctx := context.Background()

	request := &greetpb.PrimeRequest{
		Number: 120,
	}

	stream, err := c.Prime(ctx, request)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC %v", err)
	}
	defer stream.CloseSend()

LOOP:
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// we've reached the end of the stream
				break LOOP
			}
			log.Fatalf("error while reciving from GreetManyTimes RPC %v", err)
		}
		log.Printf("response from Prime:%v \n", res.GetResult())
	}

}
