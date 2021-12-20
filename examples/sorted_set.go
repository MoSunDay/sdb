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
	zpushResponse, err := c.ZPush(context.Background(),
		&pb2.ZPushRequest{Key: []byte("h"),
			Tuples: []*pb2.Tuple{
				{Value: []byte("aaa"), Score: 1.0},
				{Value: []byte("ddd"), Score: 0.8},
				{Value: []byte("bbb"), Score: 1.1},
				{Value: []byte("ccc"), Score: 0.9},
				{Value: []byte("eee"), Score: 0.7},
			}})
	log.Printf("zpushResponse: %+v, err: %+v", zpushResponse, err)
	zmembersResponse, _ := c.ZMembers(context.Background(),
		&pb2.ZMembersRequest{Key: []byte("h")})
	log.Printf("zmembersResponse: %+v, err: %+v", zmembersResponse, err)
	zrangeResponse, err := c.ZRange(context.Background(),
		&pb2.ZRangeRequest{Key: []byte("h"), Offset: 1, Limit: 100})
	log.Printf("zrangeResponse: %+v, err: %+v", zrangeResponse, err)
	zrangeResponse, err = c.ZRange(context.Background(),
		&pb2.ZRangeRequest{Key: []byte("h"), Offset: -1, Limit: 100})
	log.Printf("zrangeResponse: %+v, err: %+v", zrangeResponse, err)
	zpopResponse, err := c.ZPop(context.Background(),
		&pb2.ZPopRequest{Key: []byte("h"), Values: [][]byte{[]byte("aaa"), []byte("bbb")}})
	log.Printf("zpopResponse: %+v, err: %+v", zpopResponse, err)
	zrangeResponse, err = c.ZRange(context.Background(),
		&pb2.ZRangeRequest{Key: []byte("h"), Offset: 0, Limit: 100})
	log.Printf("zrangeResponse: %+v, err: %+v", zrangeResponse, err)
	zexistResponse, err := c.ZExist(context.Background(),
		&pb2.ZExistRequest{Key: []byte("h"),
			Values: [][]byte{[]byte("aaa"), []byte("ccc"), []byte("ddd")}})
	log.Printf("zexistResponse: %+v, err: %+v", zexistResponse, err)
}
