package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"movieexample.com/gen"
	"movieexample.com/movie/internal/controller/movie"
	metadataGateway "movieexample.com/movie/internal/gateway/metadata/grpc"
	ratingGateway "movieexample.com/movie/internal/gateway/rating/grpc"
	grpcHandler "movieexample.com/movie/internal/handler/grpc"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "Port to listen on")
	flag.Parse()

	// Register with consul
	registery, err := consul.NewRegistery("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registery.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registery.HealthCheck(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: ", err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registery.Deregister(ctx, instanceID, serviceName)

	log.Printf("Starting the movie service at port: %d\n", port)
	metadataGateway := metadataGateway.New(registery)
	ratingGateway := ratingGateway.New(registery)
	ctrl := movie.New(ratingGateway, metadataGateway)

	// New Code
	h := grpcHandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatal("failed to listen: err")
	}

	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMovieServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}
