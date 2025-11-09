package rating

import (
	"context"

	model "movieexample.com/rating/pkg"
)

// A facke rating repository.
type MockRatingRepository struct {
	returnGetRes    []*model.Rating
	returnIngestRes chan model.RatingEvent
	returnErr       error
}

// SetReturnValues set the values do you going to returns
// when te code are getting ratings.
func (m *MockRatingRepository) SetReturnGetValues(res []*model.
	Rating, err error) {
	m.returnGetRes = res
	m.returnErr = err
}

// SetReturnIngestValues set the values do you going to returns when you
// code are ingesting.
func (m *MockRatingRepository) SetReturnIngestValues(res chan model.RatingEvent, err error) {
	m.returnIngestRes = res
	m.returnErr = err
}

// Get is a facke mock of the function get only returns the values do you set before.
func (m *MockRatingRepository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]*model.Rating, error) {
	return m.returnGetRes, m.returnErr
}

// Put return the error do you put before in set values.
func (m *MockRatingRepository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return m.returnErr
}

// Ingest also returns the values do you put before in SetReturnIngestValues.
func (m *MockRatingRepository) Ingest(ctx context.Context) (chan model.RatingEvent, error) {
	return m.returnIngestRes, m.returnErr
}
