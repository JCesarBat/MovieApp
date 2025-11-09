package movietest

import (
	"movieexample.com/movie/internal/controller"
	metadatagrpc "movieexample.com/movie/internal/gateway/grpc/metadata"
	ratingrpc "movieexample.com/movie/internal/gateway/grpc/rating"
	"movieexample.com/movie/internal/handler"
	"movieexample.com/pkg/discovery"
)

func MovieTestGRPCServer(registry discovery.Registry) *handler.Handler {
	metadataGateway := metadatagrpc.New(registry)
	ratingGateway := ratingrpc.New(registry)
	ctrl := controller.New(ratingGateway, metadataGateway)
	h := handler.New(ctrl)

	return h
}
