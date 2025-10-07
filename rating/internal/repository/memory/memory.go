package memory

import (
	"context"

	"movieexample.com/rating/internal/repository"
	model "movieexample.com/rating/pkg"
)

// Repository defines a in memory reating repository
type Repository struct {
	data map[model.RecordType]map[model.RecordID][]*model.Rating
}

// New creates a new memory repository
func New() *Repository {
	return &Repository{data: map[model.RecordType]map[model.RecordID][]*model.Rating{}}
}

// Get retrives all ratings for a given record.
func (r *Repository) Get(_ context.Context, recordID model.RecordID, recordType model.RecordType) ([]*model.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}
	if ratings, ok := r.data[recordType][recordID]; !ok || len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}

	return r.data[recordType][recordID], nil
}

// Put adds a rating for a given record
func (r *Repository) Put(_ context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	if model.IsValidRecordType(recordType) {
		return repository.ErrInvalidRecordType
	}
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]*model.
			Rating{}
	}
	r.data[recordType][recordID] = append(r.data[recordType][recordID], rating)
	return nil
}
