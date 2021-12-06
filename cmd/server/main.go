package main

import (
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jsterling7/into-to-grpc/cachely"
)

// Server is the main concept that will house our cache of values
type server struct {
	data sync.Map // https://golang.org/pkg/sync/#Map
}

// Get will retrieve a value from the map and return it if found
// Get returns an error if the key is not found.
func (s *server) Get(ctx context.Context, req *cachely.GetRequest) (*cachely.GetResponse, error) {
	key := req.GetKey()
	log.Printf("looking up key %q\n", key)
	if v, ok := s.data.Load(key); ok {
		log.Printf("found key %q\n", key)
		return &cachely.GetResponse{
			Key:   key,
			Value: v.([]byte),
		}, nil
	}
	log.Printf("key not found %q\n", key)
	return nil, status.Errorf(codes.NotFound, "could not find key %s", key)
}

// Define the `Put` method.  It should update the existing value
func (s *server) Put(ctx context.Context, req *cachely.PutRequest) (*cachely.PutResponse, error) {
	key := req.GetKey()
	value := req.GetValue()
	log.Printf("setting key %s to value %s\n", key, value)

	s.data.Store(key, value)

	return &cachely.PutResponse{Key: key}, status.New(codes.OK, "").Err()
}

// Define the `Delete` method.
func (s *server) Delete(ctx context.Context, req *cachely.DeleteRequest) (*cachely.DeleteResponse, error) {
	key := req.GetKey()
	log.Printf("deleting key %s", key)

	s.data.Delete(key)

	return &cachely.DeleteResponse{Key: key}, status.New(codes.OK, "").Err()
}

func main() {
	// open a port to communicate on
	lis, err := net.Listen("tcp", ":5051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a new grpc server
	s := grpc.NewServer()

	// create our new server.
	cacheServer := server{data: sync.Map{}}

	// Register the new server by calling `cachely.RegisterCacheServer`
	cachely.RegisterCacheServer(s, &cacheServer)

	// Let the world know we are starting and where we are listening
	log.Printf("starting gRPC service on %s\n", lis.Addr())

	// start listening and responding
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
