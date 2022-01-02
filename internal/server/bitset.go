package server

import (
	"github.com/yemingfeng/sdb/internal/cluster"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type BitsetServer struct {
	pb.UnimplementedSDBServer
}

func (server *BitsetServer) BSCreate(_ context.Context, request *pb.BSCreateRequest) (*pb.BSCreateResponse, error) {
	res, err := cluster.Apply("BSCreate", request)
	return res.(*pb.BSCreateResponse), err
}

func (server *BitsetServer) BSDel(_ context.Context, request *pb.BSDelRequest) (*pb.BSDelResponse, error) {
	res, err := cluster.Apply("BSDel", request)
	return res.(*pb.BSDelResponse), err
}

func (server *BitsetServer) BSSetRange(_ context.Context, request *pb.BSSetRangeRequest) (*pb.BSSetRangeResponse, error) {
	res, err := cluster.Apply("BSSetRange", request)
	return res.(*pb.BSSetRangeResponse), err
}

func (server *BitsetServer) BSMSet(_ context.Context, request *pb.BSMSetRequest) (*pb.BSMSetResponse, error) {
	res, err := cluster.Apply("BSMSet", request)
	return res.(*pb.BSMSetResponse), err
}

func (server *BitsetServer) BSGetRange(_ context.Context, request *pb.BSGetRangeRequest) (*pb.BSGetRangeResponse, error) {
	res, err := service.BSGetRange(request.Key, request.Start, request.End)
	return &pb.BSGetRangeResponse{Values: res}, err
}

func (server *BitsetServer) BSMGet(_ context.Context, request *pb.BSMGetRequest) (*pb.BSMGetResponse, error) {
	res, err := service.BSMGet(request.Key, request.Bits)
	return &pb.BSMGetResponse{Values: res}, err
}

func (server *BitsetServer) BSCount(_ context.Context, request *pb.BSCountRequest) (*pb.BSCountResponse, error) {
	res, err := service.BSCount(request.Key)
	return &pb.BSCountResponse{Count: res}, err
}

func (server *BitsetServer) BSCountRange(_ context.Context, request *pb.BSCountRangeRequest) (*pb.BSCountRangeResponse, error) {
	res, err := service.BSCountRange(request.Key, request.Start, request.End)
	return &pb.BSCountRangeResponse{Count: res}, err
}
