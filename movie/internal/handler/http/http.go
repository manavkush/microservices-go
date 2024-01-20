package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"movieexample.com/movie/internal/controller/movie"
)

// Handler defines the movie handler
type Handler struct {
	ctrl *movie.Controller
}

// New creates a new movie HTTP Handler
func New(ctrl *movie.Controller) *Handler {
	return &Handler{ctrl}
}

func (h *Handler) GetMovieDetails(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	ctx := req.Context()
	details, err := h.ctrl.Get(ctx, id)
	if err != nil && errors.Is(err, movie.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("Response encode error: %v", err)
	}
}
