package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"movieexample.com/gen"
	"movieexample.com/pkg/config"
	"movieexample.com/pkg/discovery"
	discoveryconsul "movieexample.com/pkg/discovery/consul"
	"movieexample.com/rating/internal/controller"
	"movieexample.com/rating/internal/handler"
	"movieexample.com/rating/internal/repository/postgrers"
)

const ServiceName = "rating"

func main() {

	log.Printf("startign the rating service")
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
	repo, err := postgrers.New(cfg.GetDBConfig().ConnectionString)
	if err != nil {
		panic(err)
	}
	ctrl := controller.New(repo, nil)
	h := handler.NewGrpcHandler(ctrl)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", cfg.ServiceConfig.APIConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)

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
		srv.GracefulStop()
		log.Println("Gracefully stopped the gRPC server")
	}()	
	srv.Serve(lis)
	wg.Wait()
}

type Router struct {
	h *handler.Handler
}

func New(h *handler.Handler) *Router {
	return &Router{h: h}
}

func (r *Router) Handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "PUT":
		r.h.PutRating(w, req)

	case "GET":
		r.h.GetRatings(w, req)

	default:
		w.WriteHeader(http.StatusBadGateway)
	}
}
