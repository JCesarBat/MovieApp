package controller

import (
	"context"
	"errors"
	"log"

	"movieexample.com/metadata/pkg/model"
	"movieexample.com/movie/internal/gateway"
	moviemodel "movieexample.com/movie/pkg"
	ratingmodel "movieexample.com/rating/pkg"
)

// ErrNotFound is returned when the movie metadata is not
// found.
var ErrNotFound = errors.New("movie metadata not found")

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context,
		recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID ratingmodel.
		RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
}
type metadataGateway interface {
	Get(ctx context.Context, id string) (*model.
		Metadata, error)
}

// Controller defines a movie service controller
type Controller struct {
	ratingGateway
	metadataGateway
}

// New creates a new movie service controller
func New(rating ratingGateway, metadata metadataGateway) *Controller {
	return &Controller{
		ratingGateway:   rating,
		metadataGateway: metadata,
	}
}

// Get returns the movie details including the aggregated
// rating and movie metadata.
func (c *Controller) Get(ctx context.Context, id string) (*moviemodel.MovieDetails, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, err
	}
	rating, err := c.GetAggregatedRating(ctx, ratingmodel.RecordID(id), ratingmodel.RecordTypeMovie)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		log.Print("this movie dont have rating ")
	} else if err != nil {
		return nil, err
	}
	return &moviemodel.MovieDetails{
		Metadata: metadata,
		Rating:   &rating,
	}, nil
}
