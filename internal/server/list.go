package server

import (
	pb2 "github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type ListServer struct {
	pb2.UnimplementedSDBServer
}

func (server *ListServer) LRPush(_ context.Context, request *pb2.LRPushRequest) (*pb2.LRPushResponse, error) {
	res, err := service.LRPush(request.Key, request.Values)
	return &pb2.LRPushResponse{Success: res}, err
}

func (server *ListServer) LLPush(_ context.Context, request *pb2.LLPushRequest) (*pb2.LLPushResponse, error) {
	res, err := service.LLPush(request.Key, request.Values)
	return &pb2.LLPushResponse{Success: res}, err
}

func (server *ListServer) LPop(_ context.Context, request *pb2.LPopRequest) (*pb2.LPopResponse, error) {
	res, err := service.LPop(request.Key, request.Values)
	return &pb2.LPopResponse{Success: res}, err
}

func (server *ListServer) LRange(_ context.Context, request *pb2.LRangeRequest) (*pb2.LRangeResponse, error) {
	res, err := service.LRange(request.Key, request.Offset, request.Limit)
	return &pb2.LRangeResponse{Values: res}, err
}

func (server *ListServer) LExist(_ context.Context, request *pb2.LExistRequest) (*pb2.LExistResponse, error) {
	res, err := service.LExist(request.Key, request.Values)
	return &pb2.LExistResponse{Exists: res}, err
}

func (server *ListServer) LDel(_ context.Context, request *pb2.LDelRequest) (*pb2.LDelResponse, error) {
	res, err := service.LDel(request.Key)
	return &pb2.LDelResponse{Success: res}, err
}

func (server *ListServer) LCount(_ context.Context, request *pb2.LCountRequest) (*pb2.LCountResponse, error) {
	res, err := service.LCount(request.Key)
	return &pb2.LCountResponse{Count: res}, err
}

func (server *ListServer) LMembers(_ context.Context, request *pb2.LMembersRequest) (*pb2.LMembersResponse, error) {
	res, err := service.LMembers(request.Key)
	return &pb2.LMembersResponse{Values: res}, err
}
