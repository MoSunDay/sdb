package server

import (
	"github.com/yemingfeng/sdb/internal/cluster"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type StringServer struct {
	pb.UnimplementedSDBServer
}

func (server *StringServer) Set(_ context.Context, request *pb.SetRequest) (*pb.SetResponse, error) {
	err := cluster.Apply("Set", request)
	return &pb.SetResponse{Success: err == nil}, err
}

func (server *StringServer) MSet(_ context.Context, request *pb.MSetRequest) (*pb.MSetResponse, error) {
	err := cluster.Apply("MSet", request)
	return &pb.MSetResponse{Success: err == nil}, err
}

func (server *StringServer) SetNX(_ context.Context, request *pb.SetNXRequest) (*pb.SetNXResponse, error) {
	err := cluster.Apply("SetNX", request)
	return &pb.SetNXResponse{Success: err == nil}, err
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
	err := cluster.Apply("Del", request)
	return &pb.DelResponse{Success: err == nil}, err
}

func (server *StringServer) Incr(_ context.Context, request *pb.IncrRequest) (*pb.IncrResponse, error) {
	err := cluster.Apply("Incr", request)
	return &pb.IncrResponse{Success: err == nil}, err
}
