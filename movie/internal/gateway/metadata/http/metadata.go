package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"movieexample.com/metadata/pkg/model"
	"movieexample.com/movie/internal/gateway"
)

// Gateway defines an HTTP gateway for movie metadata service
type Gateway struct {
	addr string
}

// New creates a new HTTP gateway for movie metadata service
func New(addr string) *Gateway {
	return &Gateway{addr: addr}
}

// Get retrieves movie metadata by movie id
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, g.addr+"/metadata", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %d", resp.StatusCode)
	}
	var v *model.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
