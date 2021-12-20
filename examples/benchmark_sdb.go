package main

import (
	"context"
	pb2 "github.com/yemingfeng/sdb/internal/pb"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"strconv"
	"sync"
)

var c pb2.SDBClient = nil

func init() {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Printf("faild to connect: %+v", err)
	}
	c = pb2.NewSDBClient(conn)
}

func set(key, value []byte) {
	_, err := c.Set(context.Background(), &pb2.SetRequest{Key: key, Value: value})
	if err != nil {
		log.Fatalf("%+v, key = %s, value = %s", err, key, value)
	}
}

func get(key []byte) {
	_, err := c.Get(context.Background(), &pb2.GetRequest{Key: key})
	if err != nil {
		log.Fatalf("%+v, key = %s", err, key)
	}
}

func randBytes() []byte {
	return []byte("hello" + strconv.Itoa(rand.Int()%10000))
}

func main() {
	s := semaphore.NewWeighted(200)
	for true {
		wg := sync.WaitGroup{}
		for j := 0; j < 100000; j++ {
			wg.Add(1)
			go func() {
				_ = s.Acquire(context.Background(), 1)
				set(randBytes(), randBytes())
				wg.Done()
				s.Release(1)
			}()

			wg.Add(1)
			go func() {
				_ = s.Acquire(context.Background(), 1)
				get(randBytes())
				wg.Done()
				s.Release(1)
			}()
		}
		wg.Wait()
	}
}
