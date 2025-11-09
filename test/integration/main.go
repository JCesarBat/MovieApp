package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"movieexample.com/gen"
	"movieexample.com/metadata/pkg/metadatatest"
	"movieexample.com/movie/pkg"
	"movieexample.com/movie/pkg/movietest"
	"movieexample.com/rating/pkg/ratingtest"

	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/memory"
)

const (
	metadataServiceName = "metadata"
	ratingServiceName   = "rating"
	movieServiceName    = "movie"
	metadataServiceAddr = "localhost:8081"
	ratingServiceAddr   = "localhost:8082"
	movieServiceAddr    = "localhost:8083"
)

func main() {
	log.Println("Starting the integration test")
	ctx := context.Background()
	registry := memory.New()

	log.Println("Setting up service handlers and clients")
	metadataSrv := startMetadataService(ctx, registry)
	defer metadataSrv.GracefulStop()
	ratingSrv := startRatingService(ctx, registry)
	defer ratingSrv.GracefulStop()
	go startMovieService(ctx, registry)

	opts := grpc.WithTransportCredentials(insecure.
		NewCredentials())

	// metadata conecction client
	metadataConn, err := grpc.NewClient(metadataServiceAddr,
		opts)
	if err != nil {
		panic(err)
	}
	defer metadataConn.Close()
	metadataClient := gen.NewMetadataServiceClient(metadataConn)

	// rating conecction client
	ratingConn, err := grpc.NewClient(ratingServiceAddr,
		opts)
	if err != nil {
		panic(err)
	}
	defer metadataConn.Close()
	ratingClient := gen.NewRatingServiceClient(ratingConn)
	log.Println("Saving test metadata via metadata service")

	m := &gen.Metadata{
		Id:          "the-movie",
		Title:       "The Movie",
		Description: "The Movie, the one and only",
		Director:    "Mr. D",
	}

	if _, err := metadataClient.PutMetadata(ctx, &gen.
		PutMetadataRequest{Metadata: m}); err != nil {
		log.Fatalf("put metadata: %v", err)
	}

	log.Println("Retrieving test metadata via metadata service")
	getMetadataResp, err := metadataClient.GetMetadata(ctx,
		&gen.GetMetadataRequest{MovieId: m.Id})
	if err != nil {
		log.Fatalf("get metadata: %v", err)
	}
	if !CompareMetadataValues(getMetadataResp, m) {
		log.Println("the results are incorrect")
		return
	}
	log.Println("Getting movie details via movie service")

	// getting the movies details from a http request
	url := "http://" + movieServiceAddr + "/movie"
	log.Printf("calling the movie service with addr %v", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Panic(err)
		return
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", m.Id)
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {

		log.Panic(err)
		return
	} else if resp.StatusCode/100 != 2 {
		log.Panic(fmt.Errorf("non-2xx response : %v", resp))
		return
	}
	var v *pkg.MovieDetails
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		log.Panic(err)
		return
	}
	log.Printf("have to sey the-movie :%s", v.Metadata.ID)

	log.Println("Saving first rating via rating service")
	const userID = "user0"
	const recordTypeMovie = "movie"
	firstRating := int32(5)

	if _, err = ratingClient.PutRating(ctx, &gen.PutRatingRequest{
		UserId:     userID,
		RecordId:   m.Id,
		RecordType: recordTypeMovie,
		Value:      float32(firstRating),
	}); err != nil {
		log.Fatalf("put rating: %v", err)
	}
	log.Println("Retrieving initial aggregated rating via rating service")
	getAggregatedRatingResp, err := ratingClient.
		GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{
			RecordId:   m.Id,
			RecordType: recordTypeMovie})
	if err != nil {
		log.Fatalf("get aggreggated rating: %v", err)
	}
	if got, want := getAggregatedRatingResp.RatingValue,
		float64(5); got != want {
		log.Fatalf("rating mismatch: got %v want %v", got,
			want)
	}
	

}
func CompareMetadataValues(resp *gen.GetMetadataResponse, m *gen.Metadata) bool {
	if resp.Metadata.Id == m.Id && resp.Metadata.Title == m.Title &&
		resp.Metadata.Director == m.Director && resp.Metadata.Description == m.Description {
		return true
	}
	return false
}
func startRatingService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	log.Println("Starting rating service on" +
		ratingServiceAddr)
	h := ratingtest.RatingTestGRPCServer()
	lis, err := net.Listen("tcp", ratingServiceAddr)
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
		return nil
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()
	id := discovery.GenerateInstanceID(ratingServiceName)
	if err := registry.Register(ctx, id, ratingServiceName,
		ratingServiceAddr); err != nil {
		panic(err)
	}
	return srv
}

func startMetadataService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	log.Println("Starting metadata service on " +
		metadataServiceAddr)
	h := metadatatest.NewTestMetadataGRPCServer()
	lis, err := net.Listen("tcp", metadataServiceAddr)
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, h)
	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()
	id := discovery.GenerateInstanceID(metadataServiceName)
	if err := registry.Register(ctx, id, metadataServiceName,
		metadataServiceAddr); err != nil {
		panic(err)
	}
	return srv
}

func startMovieService(ctx context.Context, registry discovery.Registry) {
	log.Println("Starting movie service on " +
		movieServiceAddr)
	h := movietest.MovieTestGRPCServer(registry)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	instanceID := discovery.GenerateInstanceID(movieServiceName)
	if err := registry.Register(ctx, instanceID, movieServiceName, movieServiceAddr); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(movieServiceAddr, nil); err != nil {
		panic(err)
	}

}
