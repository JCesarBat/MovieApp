package pkg

import "movieexample.com/metadata/pkg/model"

//MovieDetails include movie metadata its aggregated
//rating.
type MovieDetails struct {
	Rating   *float64        `json:"rating,omitEmpty"`
	Metadata *model.Metadata `json:"metadata"`
}
