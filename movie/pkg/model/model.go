package model

import "movieexample.com/metadata/pkg/model"

// MovieDetails includes movie metadata and it's aggregated rating
type MovieDetails struct {
	Rating   *float64       `json:"rating"`
	Metadata model.Metadata `json:"meta_data"`
}
