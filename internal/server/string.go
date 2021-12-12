package server

import (
	"github.com/yemingfeng/sdb/internal/service"
	"github.com/yemingfeng/sdb/pkg/pb"
	"golang.org/x/net/context"
)

type StringServer struct {
	pb.UnimplementedSDBServer
}

func (server *StringServer) Set(_ context.Context, request *pb.SetRequest) (*pb.SetResponse, error) {
	res, err := service.Set(request.Key, request.Value)
	return &pb.SetResponse{Success: res}, err
}

func (server *StringServer) Get(_ context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	value, err := service.Get(request.Key)
	return &pb.GetResponse{Value: value}, err
}

func (server *StringServer) Del(_ context.Context, request *pb.DelRequest) (*pb.DelResponse, error) {
	res, err := service.Del(request.Key)
	return &pb.DelResponse{Success: res}, err
}

func (server *StringServer) Incr(_ context.Context, request *pb.IncrRequest) (*pb.IncrResponse, error) {
	res, err := service.Incr(request.Key, request.Delta)
	return &pb.IncrResponse{Success: res}, err
}
