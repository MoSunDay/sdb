package server

import (
	"github.com/yemingfeng/sdb/internal/service"
	"github.com/yemingfeng/sdb/pkg/pb"
	"golang.org/x/net/context"
)

type HyperLogLogServer struct {
	pb.UnimplementedSDBServer
}

func (server *HyperLogLogServer) HLLCreate(_ context.Context, request *pb.HLLCreateRequest) (*pb.HLLCreateResponse, error) {
	res, err := service.HLLCreate(request.Key)
	return &pb.HLLCreateResponse{Success: res}, err
}

func (server *HyperLogLogServer) HLLDel(_ context.Context, request *pb.HLLDelRequest) (*pb.HLLDelResponse, error) {
	res, err := service.HLLDel(request.Key)
	return &pb.HLLDelResponse{Success: res}, err
}

func (server *HyperLogLogServer) HLLAdd(_ context.Context, request *pb.HLLAddRequest) (*pb.HLLAddResponse, error) {
	res, err := service.HLLAdd(request.Key, request.Values)
	return &pb.HLLAddResponse{Success: res}, err
}

func (server *HyperLogLogServer) HLLCount(_ context.Context, request *pb.HLLCountRequest) (*pb.HLLCountResponse, error) {
	res, err := service.HLLCount(request.Key)
	return &pb.HLLCountResponse{Count: res}, err
}
