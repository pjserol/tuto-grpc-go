package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/pjserol/tuto-grpc-go/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {

	fmt.Println("Client Start!")

	opts := []grpc.DialOption{}

	if os.Getenv("Environment") == "local" {
		opts = append(opts, grpc.WithInsecure())
	} else {
		certFile := "ssl/ca.crt" // Certificate Authority trust certificate
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	cc, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//log.Printf("Created client :%f", c)

	doUnary(c)

	//doServerStreaming(c)

	//doClientStreaminng(c)

	//doBidirectionalStreaming(c)

	// should timeout
	//doUnaryWithDeadline(c, 1*time.Second)

	// should complete
	//doUnaryWithDeadline(c, 5*time.Second)
}

func doUnary(c greetpb.GreetServiceClient) {
	log.Println("Starting to do a Unary")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FistName: "Ping",
			LastName: "Pong",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet: %v", err)

	}

	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting to do a Server Streaming")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FistName: "Tic",
			LastName: "Tac",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		} else if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaminng(c greetpb.GreetServiceClient) {
	log.Println("Starting to do a Client Streaming")

	requets := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FistName: "Tic",
				LastName: "Tac",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FistName: "Ping",
				LastName: "Pong",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FistName: "Tim",
				LastName: "Tam",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	for _, req := range requets {
		log.Printf("Request: %v", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving stream: %v", err)
	}

	log.Printf("LongGreet response: %s", res.GetResult())
}

func doBidirectionalStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting to do a BiDi Streaming")

	requets := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FistName: "Tic",
				LastName: "Tac",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FistName: "Ping",
				LastName: "Pong",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FistName: "Tim",
				LastName: "Tam",
			},
		},
	}

	// create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while calling GreetEveryone: %v", err)
	}

	waitc := make(chan struct{})

	// send the messages
	go func() {
		for _, req := range requets {
			log.Printf("Sending message: %v", req)
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
		if err := stream.CloseSend(); err != nil {
			log.Fatalf("error while closing send GreetEveryone: %v", err)
		}
	}()

	// receives the messages
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("error while receiving GreetEveryone: %v", err)
				break
			}
			log.Printf("Received: %s\n", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
	log.Println("GreetEveryone finished")
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	log.Println("Starting to do a Unary with Deadline")

	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FistName: "Ping",
			LastName: "Pong",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Println("The timeout was hit! Deadline was exceeded")
			} else {
				log.Printf("Unexpected error: %v", statusErr)
			}
		} else {
			log.Fatalf("error while calling GreetWithDeadLine: %v", err)
		}
	} else {
		log.Printf("Response from GreetWithDeadLine: %v", res.Result)
	}
}
