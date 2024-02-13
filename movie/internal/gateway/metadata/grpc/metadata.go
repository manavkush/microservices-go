package grpc

import (
	"context"

	"movieexample.com/gen"
	"movieexample.com/internal/grpcutil"
	"movieexample.com/metadata/pkg/model"
	"movieexample.com/pkg/discovery"
)

// Gateway defines a gRPC gateway for movie metadata service
type Gateway struct {
	registry discovery.Registery
}

// New creates a new gRPC gateway for movie metadata service
func New(registry discovery.Registery) *Gateway {
	return &Gateway{registry: registry}
}

// Get retrieves movie metadata by movie id
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	// Create a gRPC connection to the movie metadata service
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
	if err != nil {
		return nil, err
	}
	return model.MetadataFromProto(resp.Metadata), nil
}

func (g *Gateway) Put(ctx context.Context, id string, m *model.Metadata) error {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := gen.NewMetadataServiceClient(conn)
	_, err = client.PutMetadata(ctx, &gen.PutMetadataRequest{Metadata: model.MetadataToProto(m)})
	if err != nil {
		return err
	}
	return nil
}
