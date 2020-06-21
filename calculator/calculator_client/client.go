package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/pjserol/tuto-grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {

	fmt.Println("Client Start!")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnary(c)

	//doServerStreaming(c)

	//doClientStreaming(c)

	//doBiDiStreaming(c)

	// call with a success response
	//doSquareRoot(c, 16)

	// call with an error
	//doSquareRoot(c, -5)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting to do a Unary RPC")
	req := &calculatorpb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 25,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)

	}

	log.Printf("Response from Sum: %d", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting to do a Server Streaming RPC")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 125,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		} else if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		fmt.Printf("%d ", msg.GetPrimeFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting to do a Client Streaming")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}

	numbers := []int32{12, 4, 23, 13, 20}

	for _, number := range numbers {
		fmt.Printf("Sendingn number: %d\n", number)
		if err := stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		}); err != nil {
			log.Printf("Error sending number: %s", err.Error())
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving ComputeAverage response: %v", err)
	}

	log.Printf("Response from Compute: %v", res.GetAverage())
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting to do a BiDi Streaming")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while calling FindMaximum: %v", err)
	}

	waitc := make(chan struct{})

	// Send number to the server
	go func() {
		numbers := []int32{3, 7, 2, 15, 22, 10, 8, 24, 1}

		for _, number := range numbers {
			if err := stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			}); err != nil {
				log.Printf("Error sending number: %s", err.Error())
			}
			fmt.Printf("Sending maximum: %d\n", number)
			time.Sleep(1000 * time.Millisecond)
		}
		if err := stream.CloseSend(); err != nil {
			log.Printf("Error to close send: %s", err.Error())
		}
	}()

	// Receive max from server
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("error while reading stream: %v", err)
			}

			fmt.Printf("Received maximum: %d\n", res.GetMaximum())
		}
		close(waitc)
	}()

	<-waitc

}

func doSquareRoot(c calculatorpb.CalculatorServiceClient, n int32) {
	log.Println("Starting to do a Unary - SquareRoot")

	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: n,
	})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// user error
			log.Printf("Error::Code::%d::Message::%s", respErr.Code(), respErr.Message())
			if respErr.Code() == codes.InvalidArgument {
				log.Println("Negative number not allowed!")
			}
		} else {
			// framework error
			log.Fatalf("error while calling SquareRoot: %v", err)
		}

	} else {
		fmt.Printf("SUCCESS::\nResponse from SquareRoot: %v\n", res.GetNumberRoot())
	}
}
