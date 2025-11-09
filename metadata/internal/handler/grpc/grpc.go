package grpchandler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/gen"
	"movieexample.com/metadata/internal/controller/metadata"
	"movieexample.com/metadata/internal/repository"
	"movieexample.com/metadata/pkg/model"
)

type Handler struct {
	ctrl *metadata.Controller
	gen.UnimplementedMetadataServiceServer
}

func New(ctrl *metadata.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

func (h *Handler) GetMetadata(ctx context.Context, m *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	if m == nil || m.MovieId == "" {

		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	meta, err := h.ctrl.Get(ctx, m.MovieId)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, "not found movie")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return &gen.GetMetadataResponse{
		Metadata: model.MetadataToProto(meta),
	}, nil
}
func (h *Handler) PutMetadata(ctx context.Context, m *gen.PutMetadataRequest) (*gen.PutMetadataResponse, error) {
	if m == nil || m.Metadata.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	err := h.ctrl.Put(ctx, m.Metadata.Id, model.MetadataFromProto(m.Metadata))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return &gen.PutMetadataResponse{}, nil
}
