package grpc

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/gen"
	"movieexample.com/metadata/pkg/model"
	"movieexample.com/movie/internal/controller/movie"
)

type Handler struct {
	gen.UnimplementedMovieServiceServer
	ctrl *movie.Controller
}

func New(ctrl *movie.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetMovieDetails(ctx context.Context, req *gen.GetMovieDetailsRequest) (*gen.GetMovieDetailsResponse, error) {
	log.Printf("GetMovieDetails called: %v", req.MovieId)

	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	m, err := h.ctrl.Get(ctx, req.MovieId)

	// Handle the errors returned by the grpc response
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return nil, status.Errorf(codes.NotFound, err.Error())
			default:
				log.Printf("GetMovieDetails failed: Err: %v\n", err)
				return nil, status.Errorf(codes.Internal, err.Error())
			}
		}
	}

	return &gen.GetMovieDetailsResponse{
		MovieDetails: &gen.MovieDetails{
			Rating:   *m.Rating,
			Metadata: model.MetadataToProto(&m.Metadata),
		},
	}, nil
}
