package server

import (
	pb2 "github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type StringServer struct {
	pb2.UnimplementedSDBServer
}

func (server *StringServer) Set(_ context.Context, request *pb2.SetRequest) (*pb2.SetResponse, error) {
	res, err := service.Set(request.Key, request.Value)
	return &pb2.SetResponse{Success: res}, err
}

func (server *StringServer) MSet(_ context.Context, request *pb2.MSetRequest) (*pb2.MSetResponse, error) {
	res, err := service.MSet(request.Keys, request.Values)
	return &pb2.MSetResponse{Success: res}, err
}

func (server *StringServer) SetNX(_ context.Context, request *pb2.SetNXRequest) (*pb2.SetNXResponse, error) {
	res, err := service.SetNX(request.Key, request.Value)
	return &pb2.SetNXResponse{Success: res}, err
}

func (server *StringServer) SetGet(_ context.Context, request *pb2.SetGetRequest) (*pb2.SetGetResponse, error) {
	res, old, err := service.SetGet(request.Key, request.Value)
	return &pb2.SetGetResponse{Success: res, OldValue: old}, err
}

func (server *StringServer) Get(_ context.Context, request *pb2.GetRequest) (*pb2.GetResponse, error) {
	value, err := service.Get(request.Key)
	return &pb2.GetResponse{Value: value}, err
}

func (server *StringServer) MGet(_ context.Context, request *pb2.MGetRequest) (*pb2.MGetResponse, error) {
	values, err := service.MGet(request.Keys)
	return &pb2.MGetResponse{Values: values}, err
}

func (server *StringServer) Del(_ context.Context, request *pb2.DelRequest) (*pb2.DelResponse, error) {
	res, err := service.Del(request.Key)
	return &pb2.DelResponse{Success: res}, err
}

func (server *StringServer) Incr(_ context.Context, request *pb2.IncrRequest) (*pb2.IncrResponse, error) {
	res, err := service.Incr(request.Key, request.Delta)
	return &pb2.IncrResponse{Success: res}, err
}
