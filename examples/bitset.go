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
	// 发起 bscreate 请求
	bsCreateResponse, err := c.BSCreate(context.Background(),
		&pb2.BSCreateRequest{Key: []byte("hello"), Size: 10000})
	log.Printf("bsCreateResponse: %+v, err: %+v", bsCreateResponse, err)
	// 发起 mset 请求
	bsMSetResponse, err := c.BSMSet(context.Background(),
		&pb2.BSMSetRequest{Key: []byte("hello"), Bits: []uint32{1, 2, 3}, Value: true})
	log.Printf("bsMSetResponse: %+v, err: %+v", bsMSetResponse, err)
	// 发起 mget 请求
	bsMGetResponse, err := c.BSMGet(context.Background(),
		&pb2.BSMGetRequest{Key: []byte("hello"), Bits: []uint32{4, 1, 2, 3, 5}})
	log.Printf("bsMGetResponse: %+v, err: %+v", bsMGetResponse, err)
	// 发起 setrange 请求
	bsSetResponse, err := c.BSSetRange(context.Background(),
		&pb2.BSSetRangeRequest{Key: []byte("hello"), Start: 10, End: 20, Value: true})
	log.Printf("bsSetResponse: %+v, err: %+v", bsSetResponse, err)
	// 发起 getrange 请求
	bsGetResponse, err := c.BSGetRange(context.Background(),
		&pb2.BSGetRangeRequest{Key: []byte("hello"), Start: 9, End: 21})
	log.Printf("bsGetResponse: %+v, err: %+v", bsGetResponse, err)
	// 发起 count range 请求
	bsCountRangeResponse, err := c.BSCountRange(context.Background(),
		&pb2.BSCountRangeRequest{Key: []byte("hello"), Start: 0, End: 100})
	log.Printf("bsCountRangeResponse: %+v, err: %+v", bsCountRangeResponse, err)
	// 发起 count 请求
	bsCountResponse, err := c.BSCount(context.Background(),
		&pb2.BSCountRequest{Key: []byte("hello")})
	log.Printf("bsCountResponse: %+v, err: %+v", bsCountResponse, err)
}
