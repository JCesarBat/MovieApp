package metadatatest

import (
	"movieexample.com/metadata/internal/controller/metadata"
	grpchandler "movieexample.com/metadata/internal/handler/grpc"
	"movieexample.com/metadata/internal/repository/moemory"
)

func NewTestMetadataGRPCServer() *grpchandler.Handler {
	repo := moemory.New()
	ctrl := metadata.New(repo)
	h := grpchandler.New(ctrl)
	return h
}
