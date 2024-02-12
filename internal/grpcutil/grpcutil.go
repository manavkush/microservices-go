package grpcutil

import (
	"context"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"movieexample.com/pkg/discovery"
)

// ServiceConnection attemps to select a random service instance and returns a gRPC connection to it.
func ServiceConnection(ctx context.Context, serviceName string, registery discovery.Registery) (*grpc.ClientConn, error) {
	addrs, err := registery.Discover(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	return grpc.Dial(addrs[rand.Intn(len(addrs))], grpc.WithTransportCredentials(insecure.NewCredentials()))
}
