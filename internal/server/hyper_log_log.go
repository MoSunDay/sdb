package server

import (
	"github.com/yemingfeng/sdb/internal/cluster"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type HyperLogLogServer struct {
	pb.UnimplementedSDBServer
}

func (server *HyperLogLogServer) HLLCreate(_ context.Context, request *pb.HLLCreateRequest) (*pb.HLLCreateResponse, error) {
	err := cluster.Apply("HLLCreate", request)
	return &pb.HLLCreateResponse{Success: err == nil}, err
}

func (server *HyperLogLogServer) HLLDel(_ context.Context, request *pb.HLLDelRequest) (*pb.HLLDelResponse, error) {
	err := cluster.Apply("HLLDel", request)
	return &pb.HLLDelResponse{Success: err == nil}, err
}

func (server *HyperLogLogServer) HLLAdd(_ context.Context, request *pb.HLLAddRequest) (*pb.HLLAddResponse, error) {
	err := cluster.Apply("HLLAdd", request)
	return &pb.HLLAddResponse{Success: err == nil}, err
}

func (server *HyperLogLogServer) HLLCount(_ context.Context, request *pb.HLLCountRequest) (*pb.HLLCountResponse, error) {
	res, err := service.HLLCount(request.Key)
	return &pb.HLLCountResponse{Count: res}, err
}
