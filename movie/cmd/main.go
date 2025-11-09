package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	log.Printf("startign the movie service")

	cfg := config.GetConfig()
	registry, err := discoveryconsul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
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
	s := handler.NewHttpServer(h, fmt.Sprintf("localhost:%s", cfg.ServiceConfig.APIConfig.Port))
	httpServer := handler.NewServer(s)

	// Configuring the gracefullShutDown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := <-sigChan
		cancel()
		log.Printf("Received signal %v, attempting graceful shutdown", s)
		httpServer.Shutdown(ctx)
		log.Println("Gracefully stopped the http server")
	}()
	if err := httpServer.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
	wg.Wait()
}
