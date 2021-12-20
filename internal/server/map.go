package server

import (
	pb2 "github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type MapServer struct {
	pb2.UnimplementedSDBServer
}

func (server *MapServer) MPush(_ context.Context, request *pb2.MPushRequest) (*pb2.MPushResponse, error) {
	res, err := service.MPush(request.Key, request.Pairs)
	return &pb2.MPushResponse{Success: res}, err
}

func (server *MapServer) MPop(_ context.Context, request *pb2.MPopRequest) (*pb2.MPopResponse, error) {
	res, err := service.MPop(request.Key, request.Keys)
	return &pb2.MPopResponse{Success: res}, err
}

func (server *MapServer) MExist(_ context.Context, request *pb2.MExistRequest) (*pb2.MExistResponse, error) {
	res, err := service.MExist(request.Key, request.Keys)
	return &pb2.MExistResponse{Exists: res}, err
}

func (server *MapServer) MDel(_ context.Context, request *pb2.MDelRequest) (*pb2.MDelResponse, error) {
	res, err := service.MDel(request.Key)
	return &pb2.MDelResponse{Success: res}, err
}

func (server *MapServer) MCount(_ context.Context, request *pb2.MCountRequest) (*pb2.MCountResponse, error) {
	res, err := service.MCount(request.Key)
	return &pb2.MCountResponse{Count: res}, err
}

func (server *MapServer) MMembers(_ context.Context, request *pb2.MMembersRequest) (*pb2.MMembersResponse, error) {
	res, err := service.MMembers(request.Key)
	return &pb2.MMembersResponse{Pairs: res}, err
}
