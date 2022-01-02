package server

import (
	"github.com/yemingfeng/sdb/internal/cluster"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type BloomFilterServer struct {
	pb.UnimplementedSDBServer
}

func (server *BloomFilterServer) BFCreate(_ context.Context, request *pb.BFCreateRequest) (*pb.BFCreateResponse, error) {
	res, err := cluster.Apply("BFCreate", request)
	return res.(*pb.BFCreateResponse), err
}

func (server *BloomFilterServer) BFDel(_ context.Context, request *pb.BFDelRequest) (*pb.BFDelResponse, error) {
	res, err := cluster.Apply("BFDel", request)
	return res.(*pb.BFDelResponse), err
}

func (server *BloomFilterServer) BFAdd(_ context.Context, request *pb.BFAddRequest) (*pb.BFAddResponse, error) {
	res, err := cluster.Apply("BFAdd", request)
	return res.(*pb.BFAddResponse), err
}

func (server *BloomFilterServer) BFExist(_ context.Context, request *pb.BFExistRequest) (*pb.BFExistResponse, error) {
	res, err := service.BFExist(request.Key, request.Values)
	return &pb.BFExistResponse{Exists: res}, err
}
