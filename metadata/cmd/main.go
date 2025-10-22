package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"movieexample.com/gen"
	"movieexample.com/metadata/internal/controller/metadata"
	grpchandler "movieexample.com/metadata/internal/handler/grpc"
	"movieexample.com/metadata/internal/repository/postgres"
	"movieexample.com/pkg/config"
	"movieexample.com/pkg/discovery"
	discoveryconsul "movieexample.com/pkg/discovery/consul"
)

const ServiceName = "metadata"

func main() {
	log.Printf("startign the metadata service")

	registry, err := discoveryconsul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	cfg := config.GetConfig()
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

	repo, err := postgres.New(cfg.GetDBConnectionString())
	if err != nil {
		panic(err)
	}

	ctrl := metadata.New(repo)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", cfg.ServiceConfig.APIConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, h)
	srv.Serve(lis)
}
