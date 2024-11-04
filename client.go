package main

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"sync"

	pb "pprof-example/proto"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewYourServiceClient(conn)

	wg := sync.WaitGroup{}

	for j := 0; j < 100; j++ {
		wg.Add(1)

		go func() {

			for i := 0; i < 100000000; i++ {
				ctx := context.Background()

				key := strconv.FormatInt(rand.Int63()%10000, 10)

				// Add an item
				r, err := c.AddItem(ctx, &pb.ItemRequest{Key: key, Value: "testValue"})
				if err != nil {
					log.Printf("could not add item: %v", err)
				} else {
					log.Printf("AddItem Response: %s", r.Message)
				}

				// Get the item
				r2, err := c.GetItem(ctx, &pb.ItemRequest{Key: key})
				if err != nil {
					log.Printf("could not get item: %v", err)
				} else {
					log.Printf("GetItem Response: %s, Key: %s, Value: %s", r2.Message, r2.Key, r2.Value)
				}

				// Update the item
				r3, err := c.UpdateItem(ctx, &pb.ItemRequest{Key: key, Value: "updatedValue"})
				if err != nil {
					log.Printf("could not update item: %v", err)
				} else {
					log.Printf("UpdateItem Response: %s", r3.Message)
				}

				// Get the updated item
				r4, err := c.GetItem(ctx, &pb.ItemRequest{Key: key})
				if err != nil {
					log.Printf("could not get updated item: %v", err)
				} else {
					log.Printf("GetItem Response: %s, Key: %s, Value: %s", r4.Message, r4.Key, r4.Value)
				}

				// Delete the item
				r5, err := c.DeleteItem(ctx, &pb.ItemRequest{Key: key})
				if err != nil {
					log.Printf("could not delete item: %v", err)
				} else {
					log.Printf("DeleteItem Response: %s", r5.Message)
				}

				// Try to get the deleted item
				r6, err := c.GetItem(ctx, &pb.ItemRequest{Key: key})
				if err != nil {
					log.Printf("could not get deleted item: %v", err)
				} else {
					log.Printf("GetItem Response: %s", r6.Message)
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()
}
