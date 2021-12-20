package main

import (
	pb2 "github.com/yemingfeng/sdb/internal/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Printf("faild to connect: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	// 连接服务器
	c := pb2.NewSDBClient(conn)
	setResponse, err := c.Set(context.Background(),
		&pb2.SetRequest{Key: []byte("hello"), Value: []byte("world")})
	log.Printf("setResponse: %+v, err: %+v", setResponse, err)
	getResponse, err := c.Get(context.Background(),
		&pb2.GetRequest{Key: []byte("hello")})
	log.Printf("getResponse: %+v, err: %+v", getResponse, err)
}
