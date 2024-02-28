package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"movieexample.com/gen"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
	"movieexample.com/rating/internal/controller/rating"
	grpcHandler "movieexample.com/rating/internal/handler/grpc"
	"movieexample.com/rating/internal/ingester/kafka"
	"movieexample.com/rating/internal/repository/memory"
	"movieexample.com/rating/pkg/model"
)

const serviceName = "rating"

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API Handler port")
	flag.Parse()
	fmt.Printf("Starting the movie rating service on port %d", port)

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

	repo := memory.New()
	ingester, err := kafka.NewIngester("localhost", "rating", "ratings")
	if err != nil {
		log.Fatalf("Failed to create ingester: %s\n", err)
		ingester = nil
	}
	ctrl := rating.New(repo, ingester)

	// Start the ingester in a separate go routine
	// This will keep trying to ingest the rating events from the kafka topic every 10 seconds
	go func() {
		for {
			ratingEventChannel, err := ingester.Ingest(ctx)
			if err != nil {
				log.Fatalf("Failed to ingest: %s\n", err)
			}
			for ratingEvent := range ratingEventChannel {
				// Put the rating event in the repository
				recordId := ratingEvent.RecordID
				recordType := ratingEvent.RecordType
				rating := &model.Rating{
					RecordID:   ratingEvent.RecordID,
					RecordType: ratingEvent.RecordType,
					UserID:     ratingEvent.UserID,
					Value:      ratingEvent.Value,
				}
				log.Printf("Putting rating: %v for recordId: %s recordType: %s\n", rating, recordId, recordType)
				err := ctrl.PutRating(ctx, recordId, recordType, rating)
				if err != nil {
					log.Fatalf("Failed to put rating: %s\n", err)
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()

	h := grpcHandler.New(ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	srv.Serve(lis)
}
