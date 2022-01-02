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
	res, err := cluster.Apply("Set", request)
	return res.(*pb.SetResponse), err
}

func (server *StringServer) MSet(_ context.Context, request *pb.MSetRequest) (*pb.MSetResponse, error) {
	res, err := cluster.Apply("MSet", request)
	return res.(*pb.MSetResponse), err
}

func (server *StringServer) SetNX(_ context.Context, request *pb.SetNXRequest) (*pb.SetNXResponse, error) {
	res, err := cluster.Apply("SetNX", request)
	return res.(*pb.SetNXResponse), err
}

func (server *StringServer) SetGet(_ context.Context, request *pb.SetGetRequest) (*pb.SetGetResponse, error) {
	res, err := cluster.Apply("SetGet", request)
	return res.(*pb.SetGetResponse), err
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
	res, err := cluster.Apply("Del", request)
	return res.(*pb.DelResponse), err
}

func (server *StringServer) Incr(_ context.Context, request *pb.IncrRequest) (*pb.IncrResponse, error) {
	res, err := cluster.Apply("Incr", request)
	return res.(*pb.IncrResponse), err
}
