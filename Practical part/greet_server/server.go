package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"com.grpc.nurs/greet/greetpb"
	"google.golang.org/grpc"
)

type Server struct {
	greetpb.UnimplementedCalculatorServiceServer
}


func (s *Server) Avg(stream greetpb.CalculatorService_AvgServer) error {
	fmt.Printf("Average function was invoked with a streaming request\n")
	var result int32
	var cnt int32
	cnt = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			var r float64
			r = float64(result) / float64(cnt)
			// we have finished reading the client stream
			return stream.SendAndClose(&greetpb.AvgResponse{
				Result: r,
			})

		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		number := req.GetNumber()
		result += number
		cnt++
	}
}


func (s *Server) Prime(req *greetpb.PrimeRequest, stream greetpb.CalculatorService_PrimeServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v \n", req)
	number := int(req.GetNumber())
	for i := 2; number > i; i++ {
		for number%i == 0 {
			number = number / i
			res := &greetpb.PrimeResponse{Result: int32(i)}
			if err := stream.Send(res); err != nil {
				log.Fatalf("error while sending greet many times responses: %v", err.Error())
			}
			time.Sleep(time.Second)
		}
	}
	if number > 2 {
		res := &greetpb.PrimeResponse{Result: int32(number)}
		if err := stream.Send(res); err != nil {
			log.Fatalf("error while sending greet many times responses: %v", err.Error())
		}
		time.Sleep(time.Second)
	}

	return nil
}

func main() {
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterCalculatorServiceServer(s, &Server{})
	log.Println("Server is running on port:50051")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}
