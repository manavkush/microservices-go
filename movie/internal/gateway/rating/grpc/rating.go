package grpc

import (
	"context"
	"log"

	"movieexample.com/gen"
	"movieexample.com/internal/grpcutil"
	"movieexample.com/pkg/discovery"
	"movieexample.com/rating/pkg/model"
)

// Gateway defines a gRPC gateway for rating service
type Gateway struct {
	registry discovery.Registery
}

// New creates a new gRPC gateway for rating service
func New(registry discovery.Registery) *Gateway {
	return &Gateway{registry}
}

// GetAggregatedRating returns the aggregated rating for a record or ErrNotFound if there are no ratings for it.
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	log.Println("Connecting to rating service for GetAggregatedRating")
	client := gen.NewRatingServiceClient(conn)
	resp, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{RecordId: string(recordID), RecordType: string(recordType)})
	log.Printf("GetAggregatedRating response: %v\n", resp)
	if err != nil {
		return 0, err
	}
	return resp.RatingValue, nil
}

// PutRating adds a rating for a record
func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, userID model.UserID, value model.RatingValue) error {
	log.Println("Connecting to rating service for PutRating")
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		log.Printf("Error connecting to rating service for PutRating: %v\n", err)
		return err
	}
	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)
	putRatingRequest := &gen.PutRatingRequest{
		UserId:      string(userID),
		RecordId:    string(recordID),
		RecordType:  string(recordType),
		RecordValue: int32(value),
	}

	_, err = client.PutRating(ctx, putRatingRequest)
	if err != nil {
		return err
	}
	return nil
}
