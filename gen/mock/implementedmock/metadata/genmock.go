package metadata

import (
	"context"

	"movieexample.com/metadata/pkg/model"
)

// A facke metadata repository.
type MockMetadataRepository struct {
	returnRes *model.Metadata
	returnErr error
}

// SetReturnValues set the values do you going to returns.
func (m *MockMetadataRepository) SetReturnValues(res *model.
	Metadata, err error) {
	m.returnRes = res
	m.returnErr = err
}

// Get is a facke mock of the function get only returns the values do you set before.
func (m *MockMetadataRepository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	return m.returnRes, m.returnErr
}
