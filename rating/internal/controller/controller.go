package controller

import (
	"context"
	"errors"

	"movieexample.com/rating/internal/repository"
	model "movieexample.com/rating/pkg"
)

// ErrNotFound is returned when a requested record is not
// found.
var ErrNotFound = errors.New("not Found")

type controllerInterface interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]*model.Rating, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

// Controller define a rating service controller
type Controller struct {
	repo controllerInterface
}

// New create a rating service controller
func New(repo controllerInterface) *Controller {
	return &Controller{repo: repo}
}

// GetAggregatedRating returns the aggregated rating for
// a record or errNotFound if there a ratings for it
func (c *Controller) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordID, recordType)

	if err != nil && err == repository.ErrNotFound {
		return 0, err
	}
	sum := float64(0)
	for _, r := range ratings {
		sum = float64(r.Value)
	}
	return sum / float64(len(ratings)), nil
}

// PutRating writes a rating for a given record
func (c *Controller) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {

	return c.repo.Put(ctx, recordID, recordType, rating)

}
