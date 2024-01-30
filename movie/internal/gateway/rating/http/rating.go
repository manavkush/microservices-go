package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"movieexample.com/movie/internal/gateway"
	"movieexample.com/pkg/discovery"
	"movieexample.com/rating/pkg/model"
)

// Gateway defines an HTTP gateway for movie rating service
type Gateway struct {
	registery discovery.Registery
}

// New creates a new HTTP gateway for movie rating service
func New(registery discovery.Registery) *Gateway {
	return &Gateway{registery}
}

// GetAggregatedRating returns the aggregated rating for a record.
// It returns ErrNotFound if there are no ratings for it.
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	addrs, err := g.registery.Discover(ctx, "rating")
	if err != nil {
		return 0, err
	}
	if len(addrs) == 0 {
		return 0, fmt.Errorf("no rating service instances available")
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"
	log.Printf("Calling rating service. Request: GET %s\n", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", fmt.Sprintf("%v", recordType))

	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("non-2xx response %v", resp)
	}

	var v float64
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}

	return v, nil
}

// PutRating writes a rating
func (g *Gateway) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, userId model.UserID, value model.RatingValue) error {
	addrs, err := g.registery.Discover(ctx, "rating")
	if err != nil {
		return err
	}
	if len(addrs) == 0 {
		return fmt.Errorf("no rating service instances available")
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"
	log.Printf("Calling rating service. Request: PUT %s\n", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
	if err != nil {
		return err
	}

	values := req.URL.Query()
	values.Add("id", string(recordID))
	values.Add("type", string(recordType))
	values.Add("userId", string(userId))
	values.Add("value", string(value))
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("non-2xx response: %v", err)
	}
	return nil
}
