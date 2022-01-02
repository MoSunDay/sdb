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
	res, err := cluster.Apply("HLLCreate", request)
	return res.(*pb.HLLCreateResponse), err
}

func (server *HyperLogLogServer) HLLDel(_ context.Context, request *pb.HLLDelRequest) (*pb.HLLDelResponse, error) {
	res, err := cluster.Apply("HLLDel", request)
	return res.(*pb.HLLDelResponse), err
}

func (server *HyperLogLogServer) HLLAdd(_ context.Context, request *pb.HLLAddRequest) (*pb.HLLAddResponse, error) {
	res, err := cluster.Apply("HLLAdd", request)
	return res.(*pb.HLLAddResponse), err
}

func (server *HyperLogLogServer) HLLCount(_ context.Context, request *pb.HLLCountRequest) (*pb.HLLCountResponse, error) {
	res, err := service.HLLCount(request.Key)
	return &pb.HLLCountResponse{Count: res}, err
}
