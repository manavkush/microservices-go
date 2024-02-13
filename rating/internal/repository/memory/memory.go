package memory

import (
	"context"
	"log"

	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
)

type Repository struct {
	// data is a map of map which stores all the ratings for a particular RecordType and RecordID
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

// New creates a new repository
func New() *Repository {
	return &Repository{map[model.RecordType]map[model.RecordID][]model.Rating{}}
}

// Get retrieves all the ratings for a given record
func (r *Repository) Get(ctx context.Context, recordId model.RecordID, recordType model.RecordType) ([]model.Rating, error) {

	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}

	if ratings, ok := r.data[recordType][recordId]; !ok || len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}

	log.Printf("Found ratings: %v for recordId: %v recordType: %v\n", r.data[recordType][recordId], recordId, recordType)
	return r.data[recordType][recordId], nil
}

// Put adds a rating for a given record
func (r *Repository) Put(ctx context.Context, recordId model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]model.Rating{}
	}

	r.data[recordType][recordId] = append(r.data[recordType][recordId], *rating)
	log.Printf("Added rating: %v for recordId: %v recordType: %v\n", rating, recordId, recordType)
	return nil
}
