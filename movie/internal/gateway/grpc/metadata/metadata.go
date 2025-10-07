package metadatagrpc

import (
	"context"

	"movieexample.com/gen"
	"movieexample.com/metadata/pkg/model"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/grpcutil"
)

// Gateway defines a gRPC for a movie
// metadata service.
type Gateway struct {
	registry discovery.Registry
}

// New create a gRPC gateway for a movie
// metadata service
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

// Get gets movie metadata by a movie id
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {

	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{
		MovieId: id,
	})
	if err != nil {
		return nil, err
	}
	return model.MetadataFromProto(resp.GetMetadata()), nil
}
