package server

import (
	pb2 "github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type BitsetServer struct {
	pb2.UnimplementedSDBServer
}

func (server *BitsetServer) BSCreate(_ context.Context, request *pb2.BSCreateRequest) (*pb2.BSCreateResponse, error) {
	res, err := service.BSCreate(request.Key, request.Size)
	return &pb2.BSCreateResponse{Success: res}, err
}

func (server *BitsetServer) BSDel(_ context.Context, request *pb2.BSDelRequest) (*pb2.BSDelResponse, error) {
	res, err := service.BSDel(request.Key)
	return &pb2.BSDelResponse{Success: res}, err
}

func (server *BitsetServer) BSSetRange(_ context.Context, request *pb2.BSSetRangeRequest) (*pb2.BSSetRangeResponse, error) {
	res, err := service.BSSetRange(request.Key, request.Start, request.End, request.Value)
	return &pb2.BSSetRangeResponse{Success: res}, err
}

func (server *BitsetServer) BSMSet(_ context.Context, request *pb2.BSMSetRequest) (*pb2.BSMSetResponse, error) {
	res, err := service.BSMSet(request.Key, request.Bits, request.Value)
	return &pb2.BSMSetResponse{Success: res}, err
}

func (server *BitsetServer) BSGetRange(_ context.Context, request *pb2.BSGetRangeRequest) (*pb2.BSGetRangeResponse, error) {
	res, err := service.BSGetRange(request.Key, request.Start, request.End)
	return &pb2.BSGetRangeResponse{Values: res}, err
}

func (server *BitsetServer) BSMGet(_ context.Context, request *pb2.BSMGetRequest) (*pb2.BSMGetResponse, error) {
	res, err := service.BSMGet(request.Key, request.Bits)
	return &pb2.BSMGetResponse{Values: res}, err
}

func (server *BitsetServer) BSCount(_ context.Context, request *pb2.BSCountRequest) (*pb2.BSCountResponse, error) {
	res, err := service.BSCount(request.Key)
	return &pb2.BSCountResponse{Count: res}, err
}

func (server *BitsetServer) BSCountRange(_ context.Context, request *pb2.BSCountRangeRequest) (*pb2.BSCountRangeResponse, error) {
	res, err := service.BSCountRange(request.Key, request.Start, request.End)
	return &pb2.BSCountRangeResponse{Count: res}, err
}
