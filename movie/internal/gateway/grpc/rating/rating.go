package ratingrpc

import (
	"context"

	"movieexample.com/gen"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/grpcutil"
	model "movieexample.com/rating/pkg"
)

// Gateway defines a gRPC for a movie
// rating service.
type Gateway struct {
	registry discovery.Registry
}

// New create a gRPC gateway for a movie
// rating service
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

// GetAggregatedRating returns the aggregated rating for a
// record or ErrNotFound if there are no rating for it
func (g *Gateway) GetAggregatedRating(ctx context.Context,
	recordID model.RecordID, recordType model.RecordType) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	client := gen.NewRatingServiceClient(conn)

	req := &gen.GetAggregatedRatingRequest{
		RecordId:   string(recordID),
		RecordType: string(recordType),
	}

	resp, err := client.GetAggregatedRating(ctx, req)
	if err != nil {
		return 0, nil
	}
	return resp.RatingValue, nil
}

// PutRating writes a rating
func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID,
	recordType model.RecordType, rating *model.Rating) error {

	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewRatingServiceClient(conn)

	req := &gen.PutRatingRequest{
		RecordId:   string(recordID),
		RecordType: string(recordType),
		UserId:     string(rating.UserID),
		Value:      float32(rating.Value),
	}
	_, err = client.PutRating(ctx, req)
	return err
}
