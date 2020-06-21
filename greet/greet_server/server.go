package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/pjserol/tuto-grpc-go/greet/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("Great function invoked with %v \n", req)
	res := fmt.Sprintf("Hello, %s %s", req.GetGreeting().FistName, req.GetGreeting().LastName)
	return &greetpb.GreetResponse{
		Result: res,
	}, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Printf("GreatManyTimes function invoked with %v \n", req)

	for i := 0; i < 5; i++ {
		stream.Send(
			&greetpb.GreetManyTimesResponse{
				Result: fmt.Sprintf("Hello, %s %s number %d", req.GetGreeting().FistName, req.GetGreeting().LastName, i),
			},
		)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	log.Printf("LongGreet function invoked with a streaming request")
	result := ""

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// we've reached the end of the client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		} else if err != nil {
			return err
		}
		res := fmt.Sprintf("Hello, %s %s!", msg.GetGreeting().FistName, msg.GetGreeting().LastName)
		result += res + "\n"
	}
}

// Bi-Directional Streaming
func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	log.Printf("GreetEveryone function invoked with a streaming request")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		log.Printf("Received: %v", req)
		if err := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: fmt.Sprintf("Hello, %s %s", req.GetGreeting().GetFistName(), req.GetGreeting().GetLastName()),
		}); err != nil {
			return err
		}
	}
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	log.Printf("Great GreetWithDeadline invoked with %v \n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			// the client canceled the request
			log.Println("The client cancelled the request")
			return nil, status.Error(codes.Canceled, "The client canclled the request!")
		}
		time.Sleep(1 * time.Second)
	}
	res := fmt.Sprintf("Hello, %s %s", req.GetGreeting().FistName, req.GetGreeting().LastName)
	return &greetpb.GreetWithDeadlineResponse{
		Result: res,
	}, nil
}

func main() {
	fmt.Println("Start server!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	if os.Getenv("Environment") != "local" {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("Failed loadind certification: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v ", err)
	}
}
