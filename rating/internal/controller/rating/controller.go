package rating

import (
	"context"
	"errors"
	"log"

	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
)

// ErrNotFound is returned when no ratings are found for a record
var ErrNotFound = errors.New("ratings not found for the record")

type ratingRepository interface {
	Get(ctx context.Context, recordId model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordId model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

type ratingIngester interface {
	Ingest(ctx context.Context) (<-chan model.RatingEvent, error)
}

// Controller defines a rating service controller
type Controller struct {
	repo     ratingRepository
	ingester ratingIngester
}

// New creates a rating service controller
func New(repo ratingRepository, ingester ratingIngester) *Controller {
	return &Controller{repo, ingester}
}

// GetAggregatedRating returns the aggregated rating for a
// record or ErrNotFound if there are no ratings for it
func (c *Controller) GetAggregatedRating(ctx context.Context, recordId model.RecordID, recordType model.RecordType) (float64, error) {
	log.Printf("GetAggregatedRating called for recordId: %v recordType: %v\n", recordId, recordType)
	ratings, err := c.repo.Get(ctx, recordId, recordType)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		log.Printf("No ratings found for recordId: %v recordType: %v\n", recordId, recordType)
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	sum := float64(0)

	for _, record := range ratings {
		sum += float64(record.Value)
	}

	log.Printf("Aggregated rating: %v for recordId: %v recordType: %v\n", sum/float64(len(ratings)), recordId, recordType)
	return sum / (float64(len(ratings))), nil
}

func (c *Controller) PutRating(ctx context.Context, recordId model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return c.repo.Put(ctx, recordId, recordType, rating)
}
