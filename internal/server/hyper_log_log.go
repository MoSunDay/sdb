package server

import (
	pb2 "github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type HyperLogLogServer struct {
	pb2.UnimplementedSDBServer
}

func (server *HyperLogLogServer) HLLCreate(_ context.Context, request *pb2.HLLCreateRequest) (*pb2.HLLCreateResponse, error) {
	res, err := service.HLLCreate(request.Key)
	return &pb2.HLLCreateResponse{Success: res}, err
}

func (server *HyperLogLogServer) HLLDel(_ context.Context, request *pb2.HLLDelRequest) (*pb2.HLLDelResponse, error) {
	res, err := service.HLLDel(request.Key)
	return &pb2.HLLDelResponse{Success: res}, err
}

func (server *HyperLogLogServer) HLLAdd(_ context.Context, request *pb2.HLLAddRequest) (*pb2.HLLAddResponse, error) {
	res, err := service.HLLAdd(request.Key, request.Values)
	return &pb2.HLLAddResponse{Success: res}, err
}

func (server *HyperLogLogServer) HLLCount(_ context.Context, request *pb2.HLLCountRequest) (*pb2.HLLCountResponse, error) {
	res, err := service.HLLCount(request.Key)
	return &pb2.HLLCountResponse{Count: res}, err
}
