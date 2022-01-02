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
	res, err := cluster.Apply("SPush", request)
	return res.(*pb.SPushResponse), err
}

func (server *SetServer) SPop(_ context.Context, request *pb.SPopRequest) (*pb.SPopResponse, error) {
	res, err := cluster.Apply("SPop", request)
	return res.(*pb.SPopResponse), err
}

func (server *SetServer) SExist(_ context.Context, request *pb.SExistRequest) (*pb.SExistResponse, error) {
	res, err := service.SExist(request.Key, request.Values)
	return &pb.SExistResponse{Exists: res}, err
}

func (server *SetServer) SDel(_ context.Context, request *pb.SDelRequest) (*pb.SDelResponse, error) {
	res, err := cluster.Apply("SDel", request)
	return res.(*pb.SDelResponse), err
}

func (server *SetServer) SCount(_ context.Context, request *pb.SCountRequest) (*pb.SCountResponse, error) {
	res, err := service.SCount(request.Key)
	return &pb.SCountResponse{Count: res}, err
}

func (server *SetServer) SMembers(_ context.Context, request *pb.SMembersRequest) (*pb.SMembersResponse, error) {
	res, err := service.SMembers(request.Key)
	return &pb.SMembersResponse{Values: res}, err
}
