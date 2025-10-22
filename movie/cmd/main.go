package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"movieexample.com/movie/internal/controller"
	metadatagrpc "movieexample.com/movie/internal/gateway/grpc/metadata"
	ratingrpc "movieexample.com/movie/internal/gateway/grpc/rating"
	"movieexample.com/movie/internal/handler"
	"movieexample.com/pkg/config"
	"movieexample.com/pkg/discovery"
	discoveryconsul "movieexample.com/pkg/discovery/consul"
)

const ServiceName = "movie"

func main() {

	log.Printf("startign the metadata service")

	cfg := config.GetConfig()
	registry, err := discoveryconsul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(ServiceName)
	if err := registry.Register(ctx, instanceID, ServiceName, fmt.Sprintf("localhost:%s", cfg.ServiceConfig.APIConfig.Port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, ServiceName); err != nil {
				log.Println("Failed to report healthy state :" + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, ServiceName)
	metadataGateway := metadatagrpc.New(registry)
	ratingGateway := ratingrpc.New(registry)

	ctrl := controller.New(ratingGateway, metadataGateway)
	h := handler.New(ctrl)
	log.Printf("the server is running in port %s", cfg.ServiceConfig.APIConfig.Port)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%s", cfg.ServiceConfig.APIConfig.Port), nil); err != nil {
		panic(err)
	}
}
