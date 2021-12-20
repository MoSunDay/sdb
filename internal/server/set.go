package server

import (
	pb2 "github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type SetServer struct {
	pb2.UnimplementedSDBServer
}

func (server *SetServer) SPush(_ context.Context, request *pb2.SPushRequest) (*pb2.SPushResponse, error) {
	res, err := service.SPush(request.Key, request.Values)
	return &pb2.SPushResponse{Success: res}, err
}

func (server *SetServer) SPop(_ context.Context, request *pb2.SPopRequest) (*pb2.SPopResponse, error) {
	res, err := service.SPop(request.Key, request.Values)
	return &pb2.SPopResponse{Success: res}, err
}

func (server *SetServer) SExist(_ context.Context, request *pb2.SExistRequest) (*pb2.SExistResponse, error) {
	res, err := service.SExist(request.Key, request.Values)
	return &pb2.SExistResponse{Exists: res}, err
}

func (server *SetServer) SDel(_ context.Context, request *pb2.SDelRequest) (*pb2.SDelResponse, error) {
	res, err := service.SDel(request.Key)
	return &pb2.SDelResponse{Success: res}, err
}

func (server *SetServer) SCount(_ context.Context, request *pb2.SCountRequest) (*pb2.SCountResponse, error) {
	res, err := service.SCount(request.Key)
	return &pb2.SCountResponse{Count: res}, err
}

func (server *SetServer) SMembers(_ context.Context, request *pb2.SMembersRequest) (*pb2.SMembersResponse, error) {
	res, err := service.SMembers(request.Key)
	return &pb2.SMembersResponse{Values: res}, err
}
