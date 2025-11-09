package ratingtest

import (
	"movieexample.com/rating/internal/controller"
	"movieexample.com/rating/internal/handler"
	"movieexample.com/rating/internal/repository/memory"
)

func RatingTestGRPCServer() *handler.GrpcHandler {
	repo := memory.New()
	ctrl := controller.New(repo, nil)
	h := handler.NewGrpcHandler(ctrl)
	return h
}
