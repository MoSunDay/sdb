package server

import (
	"github.com/yemingfeng/sdb/internal/service"
	"github.com/yemingfeng/sdb/pkg/pb"
	"golang.org/x/net/context"
)

type SortedSetServer struct {
	pb.UnimplementedSDBServer
}

func (server *SortedSetServer) ZPush(_ context.Context, request *pb.ZPushRequest) (*pb.ZPushResponse, error) {
	res, err := service.ZPush(request.Key, request.Tuples)
	return &pb.ZPushResponse{Success: res}, err
}

func (server *SortedSetServer) ZPop(_ context.Context, request *pb.ZPopRequest) (*pb.ZPopResponse, error) {
	res, err := service.ZPop(request.Key, request.Values)
	return &pb.ZPopResponse{Success: res}, err
}

func (server *SortedSetServer) ZRange(_ context.Context, request *pb.ZRangeRequest) (*pb.ZRangeResponse, error) {
	res, err := service.ZRange(request.Key, request.Offset, request.Limit)
	return &pb.ZRangeResponse{Tuples: res}, err
}

func (server *SortedSetServer) ZExist(_ context.Context, request *pb.ZExistRequest) (*pb.ZExistResponse, error) {
	res, err := service.ZExist(request.Key, request.Values)
	return &pb.ZExistResponse{Exists: res}, err
}

func (server *SortedSetServer) ZDel(_ context.Context, request *pb.ZDelRequest) (*pb.ZDelResponse, error) {
	res, err := service.ZDel(request.Key)
	return &pb.ZDelResponse{Success: res}, err
}

func (server *SetServer) ZCount(_ context.Context, request *pb.ZCountRequest) (*pb.ZCountResponse, error) {
	res, err := service.ZCount(request.Key)
	return &pb.ZCountResponse{Count: res}, err
}

func (server *SetServer) ZMembers(_ context.Context, request *pb.ZMembersRequest) (*pb.ZMembersResponse, error) {
	res, err := service.ZMembers(request.Key)
	return &pb.ZMembersResponse{Tuples: res}, err
}
