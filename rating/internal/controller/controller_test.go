package controller

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"movieexample.com/gen/mock/implementedmock/rating"
	"movieexample.com/pkg/random"
	model "movieexample.com/rating/pkg"
)

func TestGetAggregatedRating(t *testing.T) {
	var ratings []*model.Rating

	for i := 0; i < 10; i++ {
		ratings = append(ratings, random.RandomRating(t))
	}
	tests := []struct {
		name       string
		expRepoRes []*model.Rating
		expRepoErr error

		wantRes []*model.Rating
		wantErr error
	}{
		{
			name:       "success",
			expRepoRes: ratings,
			wantRes:    ratings,
		},
		{
			name:       "Not Found",
			expRepoErr: ErrNotFound,
			wantErr:    ErrNotFound,
		},
		{
			name:       "Another Error",
			expRepoErr: errors.New("New error"),
			wantErr:    errors.New("New error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := rating.MockRatingRepository{}
			repo.SetReturnGetValues(tt.expRepoRes, tt.wantErr)
			c := New(&repo, nil)
			resp, err := c.GetAggregatedRating(t.Context(), model.RecordID("some-id"), model.RecordType("some-recordType"))

			require.Equal(t, resp, CalculateValue(tt.expRepoRes))
			require.Equal(t, tt.expRepoErr, err)
		})
	}
}

func CalculateValue(ratings []*model.Rating) float64 {

	sum := float64(0)
	for _, r := range ratings {
		sum = float64(r.Value)
	}
	if sum == 0 {
		return 0
	}
	return sum / float64(len(ratings))
}
