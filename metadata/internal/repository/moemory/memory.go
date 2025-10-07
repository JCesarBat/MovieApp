package moemory

import (
	"context"
	"sync"

	"movieexample.com/metadata/internal/repository"
	"movieexample.com/metadata/pkg/model"
)

// Repository defines a memory movie metadata repository
type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

// New create a new memory repository
func New() *Repository {

	return &Repository{data: map[string]*model.Metadata{}}
}

// Get retrive movie metadata for my movie id
func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {

	r.RLock()
	defer r.RUnlock()

	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil

}

// Put adds  movie metadata for a given id
func (r *Repository) Put(_ context.Context, id string, m *model.Metadata) error {
	r.Lock()
	defer r.Unlock()

	r.data[id] = m

	return nil
}
