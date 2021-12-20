package server

import (
	pb2 "github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type BloomFilterServer struct {
	pb2.UnimplementedSDBServer
}

func (server *BloomFilterServer) BFCreate(_ context.Context, request *pb2.BFCreateRequest) (*pb2.BFCreateResponse, error) {
	res, err := service.BFCreate(request.Key, request.N, request.P)
	return &pb2.BFCreateResponse{Success: res}, err
}

func (server *BloomFilterServer) BFDel(_ context.Context, request *pb2.BFDelRequest) (*pb2.BFDelResponse, error) {
	res, err := service.BFDel(request.Key)
	return &pb2.BFDelResponse{Success: res}, err
}

func (server *BloomFilterServer) BFAdd(_ context.Context, request *pb2.BFAddRequest) (*pb2.BFAddResponse, error) {
	res, err := service.BFAdd(request.Key, request.Values)
	return &pb2.BFAddResponse{Success: res}, err
}

func (server *BloomFilterServer) BFExist(_ context.Context, request *pb2.BFExistRequest) (*pb2.BFExistResponse, error) {
	res, err := service.BFExist(request.Key, request.Values)
	return &pb2.BFExistResponse{Exists: res}, err
}
