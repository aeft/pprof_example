package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"sync"

	pb "pprof-example/proto"

	"google.golang.org/grpc"

	_ "net/http/pprof" // Enable pprof
)

const (
	port = ":50051"
)

// Global map
var dataStore sync.Map

// server is used to implement the gRPC service.
type server struct {
	pb.UnimplementedYourServiceServer
}

func magic(s string) string {
	for i := 0; i < 100; i++ {
		s += strconv.FormatInt(rand.Int63()%10, 10)
	}

	hash := sha256.New()

	hash.Write([]byte(s))

	res := []byte(s)
	for i := 0; i < 100; i++ {
		hash := sha256.Sum256(res)
		res = hash[:]
	}

	return hex.EncodeToString(res)
}

// AddItem adds a new item to the global map
func (s *server) AddItem(ctx context.Context, req *pb.ItemRequest) (*pb.OperationResponse, error) {
	dataStore.Store(req.Key, magic(req.Value))
	return &pb.OperationResponse{Message: "Item added successfully"}, nil
}

// DeleteItem deletes an item from the global map
func (s *server) DeleteItem(ctx context.Context, req *pb.ItemRequest) (*pb.OperationResponse, error) {
	_, loaded := dataStore.LoadAndDelete(req.Key)

	if loaded {
		return &pb.OperationResponse{Message: "Item deleted"}, nil
	} else {
		return &pb.OperationResponse{Message: "Item not found"}, nil
	}
}

// UpdateItem updates an existing item in the global map
func (s *server) UpdateItem(ctx context.Context, req *pb.ItemRequest) (*pb.OperationResponse, error) {
	dataStore.Store(req.Key, magic(req.Value))
	return &pb.OperationResponse{Message: "Item updated successfully"}, nil
}

// GetItem retrieves an item from the global map
func (s *server) GetItem(ctx context.Context, req *pb.ItemRequest) (*pb.ItemResponse, error) {
	if value, exists := dataStore.Load(req.Key); exists {
		return &pb.ItemResponse{Key: req.Key, Value: value.(string), Message: "Item retrieved successfully"}, nil
	}
	return &pb.ItemResponse{Message: "Item not found"}, nil
}

func main() {
	// Start pprof in a separate goroutine
	go func() {
		runtime.SetCPUProfileRate(1000)
		log.Println("Starting pprof on :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterYourServiceServer(s, &server{})

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
