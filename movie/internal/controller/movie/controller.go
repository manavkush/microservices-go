package movie

import (
	"context"
	"errors"
	"log"

	metadataModel "movieexample.com/metadata/pkg/model"
	"movieexample.com/movie/internal/gateway"
	"movieexample.com/movie/pkg/model"
	ratingModel "movieexample.com/rating/pkg/model"
)

// ErrNotFound is returned when the movie metadata is not found
var ErrNotFound = errors.New("movie metadata not found")

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID ratingModel.RecordID, recordType ratingModel.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID ratingModel.RecordID, recordType ratingModel.RecordType, userId ratingModel.UserID, value ratingModel.RatingValue) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadataModel.Metadata, error)
}

// Controller defines a movie service controller
type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

// New creates a new movie service controller
func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{
		ratingGateway,
		metadataGateway,
	}
}

// Get returns the movie details including the aggregated rating and the movie metadata.
func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	log.Printf("Movie Controller: Get called: %v", id)
	metadata, err := c.metadataGateway.Get(ctx, id)

	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		log.Printf("Movie Controller: Not found: %v", err)
		return nil, ErrNotFound
	} else if err != nil {
		log.Printf("Movie Controller: Internal Error: %v", err)
		return nil, err
	}
	log.Printf("Movie Controller: Found the movie: %v", metadata)

	details := &model.MovieDetails{Metadata: *metadata}
	rating, err := c.ratingGateway.GetAggregatedRating(ctx, ratingModel.RecordID(id), ratingModel.RecordMovieType)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	details.Rating = &rating

	log.Printf("Movie Controller: Rating: %v", rating)

	return details, nil
}
