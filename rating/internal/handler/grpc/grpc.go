package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/gen"
	"movieexample.com/rating/internal/controller/rating"
	"movieexample.com/rating/pkg/model"
)

// Handler defines a gRPC rating API handler
type Handler struct {
	gen.UnimplementedRatingServiceServer
	svc *rating.Controller
}

// New creates a new gRPC rating API handler
func New(ctrl *rating.Controller) *Handler {
	return &Handler{svc: ctrl}
}

// GetAggregateRating returns the aggregated rating for a record.
func (h *Handler) GetAggregateRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req.RecordId == "" || req.RecordType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty record id or record type")
	}
	val, err := h.svc.GetAggregatedRatings(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType))
	if err != nil && errors.Is(err, rating.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetAggregatedRatingResponse{RatingValue: val}, nil
}

// PutRating writes a rating for a given record.
func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty user id or record id")
	}
	rating := &model.Rating{
		RecordID:   model.RecordID(req.RecordId),
		RecordType: model.RecordType(req.RecordType),
		UserID:     model.UserID(req.UserId),
		Value:      model.RatingValue(req.RecordValue),
	}
	if err := h.svc.PutRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType), rating); err != nil {
		return nil, err
	}
	return &gen.PutRatingResponse{}, nil
}
