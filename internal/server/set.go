package server

import (
	"github.com/yemingfeng/sdb/internal/cluster"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type SetServer struct {
	pb.UnimplementedSDBServer
}

func (server *SetServer) SPush(_ context.Context, request *pb.SPushRequest) (*pb.SPushResponse, error) {
	err := cluster.Apply("SPush", request)
	return &pb.SPushResponse{Success: err == nil}, err
}

func (server *SetServer) SPop(_ context.Context, request *pb.SPopRequest) (*pb.SPopResponse, error) {
	err := cluster.Apply("SPop", request)
	return &pb.SPopResponse{Success: err == nil}, err
}

func (server *SetServer) SExist(_ context.Context, request *pb.SExistRequest) (*pb.SExistResponse, error) {
	res, err := service.SExist(request.Key, request.Values)
	return &pb.SExistResponse{Exists: res}, err
}

func (server *SetServer) SDel(_ context.Context, request *pb.SDelRequest) (*pb.SDelResponse, error) {
	err := cluster.Apply("SDel", request)
	return &pb.SDelResponse{Success: err == nil}, err
}

func (server *SetServer) SCount(_ context.Context, request *pb.SCountRequest) (*pb.SCountResponse, error) {
	res, err := service.SCount(request.Key)
	return &pb.SCountResponse{Count: res}, err
}

func (server *SetServer) SMembers(_ context.Context, request *pb.SMembersRequest) (*pb.SMembersResponse, error) {
	res, err := service.SMembers(request.Key)
	return &pb.SMembersResponse{Values: res}, err
}
