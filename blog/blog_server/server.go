package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/pjserol/tuto-grpc-go/blog/blogpb"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var collection *mongo.Collection

type server struct{}

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Connection to mongoDB")
	// connect to mongoDB
	// with secure mode: mongodb://foo:bar@localhost:27017
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = mongoClient.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	collection = mongoClient.Database("blogdb").Collection("blog")

	log.Println("Start blog service!")

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
	blogpb.RegisterBlogServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		log.Println("Starting blog server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v ", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	log.Println("Stopping blog server.")
	s.Stop()
	log.Println("Closing the listener of blog server.")
	lis.Close()
	log.Println("Closing mongoDB connection")
	mongoClient.Disconnect(context.TODO())
	log.Println("End of blog service.")
}
