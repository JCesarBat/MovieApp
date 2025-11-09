package random

import (
	"math/rand"
	"testing"

	model "movieexample.com/rating/pkg"
)

func RandomRating(t *testing.T) *model.Rating {

	return &model.Rating{
		RecordID:   model.RecordID("some-recordID"),
		RecordType: model.RecordType("some-recordType"),
		UserID:     model.UserID("some-userID"),
		Value:      model.RatingValue(rand.Intn(100)),
	}
}
