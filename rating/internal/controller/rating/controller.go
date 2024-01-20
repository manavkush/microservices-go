package rating

import (
	"context"
	"errors"

	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
)

// ErrNotFound is returned when no ratings are found for a record
var ErrNotFound = errors.New("ratings not found for the record")

type ratingRepository interface {
	Get(ctx context.Context, recordId model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordId model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

// Controller defines a rating service controller
type Controller struct {
	repo ratingRepository
}

// New creates a rating service controller
func New(repo ratingRepository) *Controller {
	return &Controller{repo}
}

// GetAggregatedRatings returns the aggregated rating for a
// record or ErrNotFound if there are no ratings for it
func (c *Controller) GetAggregatedRatings(ctx context.Context, recordId model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordId, recordType)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	sum := float64(0)

	for _, record := range ratings {
		sum += float64(record.Value)
	}

	return sum / (float64(len(ratings))), nil
}

func (c *Controller) PutRating(ctx context.Context, recordId model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return c.repo.Put(ctx, recordId, recordType, rating)
}
