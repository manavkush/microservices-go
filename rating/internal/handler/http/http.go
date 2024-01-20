package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"movieexample.com/rating/internal/controller/rating"
	"movieexample.com/rating/pkg/model"
)

// Handler defines a rating service controller
type Handler struct {
	ctrl *rating.Controller
}

// New creates a new rating service HTTP handler
func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl}
}

// Handle handles GET and PUT /rating requests
func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	recordID := model.RecordID(req.FormValue("id"))
	if recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recordType := model.RecordType(req.FormValue("type"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	switch req.Method {
	case http.MethodGet:
		v, err := h.ctrl.GetAggregatedRatings(ctx, recordID, recordType)
		if err != nil && errors.Is(err, rating.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("Response encode error: %v\n", err)
		}

	case http.MethodPut:
		userID := model.UserID(req.FormValue("userId"))
		v, err := strconv.ParseInt(req.FormValue("value"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ratingRecord := &model.Rating{RecordID: recordID, RecordType: recordType, UserID: userID, Value: model.RatingValue(v)}

		if err := h.ctrl.PutRating(ctx, recordID, recordType, ratingRecord); err != nil {
			log.Printf("Repository put error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
