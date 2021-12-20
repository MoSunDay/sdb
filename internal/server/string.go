package server

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type StringServer struct {
	pb.UnimplementedSDBServer
}

func (server *StringServer) Set(_ context.Context, request *pb.SetRequest) (*pb.SetResponse, error) {
	res, err := service.Set(request.Key, request.Value)
	return &pb.SetResponse{Success: res}, err
}

func (server *StringServer) MSet(_ context.Context, request *pb.MSetRequest) (*pb.MSetResponse, error) {
	res, err := service.MSet(request.Keys, request.Values)
	return &pb.MSetResponse{Success: res}, err
}

func (server *StringServer) SetNX(_ context.Context, request *pb.SetNXRequest) (*pb.SetNXResponse, error) {
	res, err := service.SetNX(request.Key, request.Value)
	return &pb.SetNXResponse{Success: res}, err
}

func (server *StringServer) SetGet(_ context.Context, request *pb.SetGetRequest) (*pb.SetGetResponse, error) {
	res, old, err := service.SetGet(request.Key, request.Value)
	return &pb.SetGetResponse{Success: res, OldValue: old}, err
}

func (server *StringServer) Get(_ context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	value, err := service.Get(request.Key)
	return &pb.GetResponse{Value: value}, err
}

func (server *StringServer) MGet(_ context.Context, request *pb.MGetRequest) (*pb.MGetResponse, error) {
	values, err := service.MGet(request.Keys)
	return &pb.MGetResponse{Values: values}, err
}

func (server *StringServer) Del(_ context.Context, request *pb.DelRequest) (*pb.DelResponse, error) {
	res, err := service.Del(request.Key)
	return &pb.DelResponse{Success: res}, err
}

func (server *StringServer) Incr(_ context.Context, request *pb.IncrRequest) (*pb.IncrResponse, error) {
	res, err := service.Incr(request.Key, request.Delta)
	return &pb.IncrResponse{Success: res}, err
}
