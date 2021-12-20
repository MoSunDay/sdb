package server

import (
	pb2 "github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type SortedSetServer struct {
	pb2.UnimplementedSDBServer
}

func (server *SortedSetServer) ZPush(_ context.Context, request *pb2.ZPushRequest) (*pb2.ZPushResponse, error) {
	res, err := service.ZPush(request.Key, request.Tuples)
	return &pb2.ZPushResponse{Success: res}, err
}

func (server *SortedSetServer) ZPop(_ context.Context, request *pb2.ZPopRequest) (*pb2.ZPopResponse, error) {
	res, err := service.ZPop(request.Key, request.Values)
	return &pb2.ZPopResponse{Success: res}, err
}

func (server *SortedSetServer) ZRange(_ context.Context, request *pb2.ZRangeRequest) (*pb2.ZRangeResponse, error) {
	res, err := service.ZRange(request.Key, request.Offset, request.Limit)
	return &pb2.ZRangeResponse{Tuples: res}, err
}

func (server *SortedSetServer) ZExist(_ context.Context, request *pb2.ZExistRequest) (*pb2.ZExistResponse, error) {
	res, err := service.ZExist(request.Key, request.Values)
	return &pb2.ZExistResponse{Exists: res}, err
}

func (server *SortedSetServer) ZDel(_ context.Context, request *pb2.ZDelRequest) (*pb2.ZDelResponse, error) {
	res, err := service.ZDel(request.Key)
	return &pb2.ZDelResponse{Success: res}, err
}

func (server *SetServer) ZCount(_ context.Context, request *pb2.ZCountRequest) (*pb2.ZCountResponse, error) {
	res, err := service.ZCount(request.Key)
	return &pb2.ZCountResponse{Count: res}, err
}

func (server *SetServer) ZMembers(_ context.Context, request *pb2.ZMembersRequest) (*pb2.ZMembersResponse, error) {
	res, err := service.ZMembers(request.Key)
	return &pb2.ZMembersResponse{Tuples: res}, err
}
