package server

import (
	"github.com/yemingfeng/sdb/internal/cluster"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type GeoHashServer struct {
	pb.UnimplementedSDBServer
}

func (server *GeoHashServer) GHCreate(_ context.Context, request *pb.GHCreateRequest) (*pb.GHCreateResponse, error) {
	err := cluster.Apply("GHCreate", request)
	return &pb.GHCreateResponse{Success: err == nil}, err
}

func (server *GeoHashServer) GHDel(_ context.Context, request *pb.GHDelRequest) (*pb.GHDelResponse, error) {
	err := cluster.Apply("GHDel", request)
	return &pb.GHDelResponse{Success: err == nil}, err
}

func (server *GeoHashServer) GHAdd(_ context.Context, request *pb.GHAddRequest) (*pb.GHAddResponse, error) {
	err := cluster.Apply("GHAdd", request)
	return &pb.GHAddResponse{Success: err == nil}, err
}

func (server *GeoHashServer) GHRem(_ context.Context, request *pb.GHRemRequest) (*pb.GHRemResponse, error) {
	err := cluster.Apply("GHRem", request)
	return &pb.GHRemResponse{Success: err == nil}, err
}

func (server *GeoHashServer) GHGetBoxes(_ context.Context, request *pb.GHGetBoxesRequest) (*pb.GHGetBoxesResponse, error) {
	res, err := service.GHGetBoxes(request.Key, request.Latitude, request.Longitude)
	return &pb.GHGetBoxesResponse{Points: res}, err
}

func (server *GeoHashServer) GHGetNeighbors(_ context.Context, request *pb.GHGetNeighborsRequest) (*pb.GHGetNeighborsResponse, error) {
	res, err := service.GHGetNeighbors(request.Key, request.Latitude, request.Longitude)
	return &pb.GHGetNeighborsResponse{Points: res}, err
}

func (server *GeoHashServer) GHCount(_ context.Context, request *pb.GHCountRequest) (*pb.GHCountResponse, error) {
	res, err := service.GHCount(request.Key)
	return &pb.GHCountResponse{Count: res}, err
}
func (server *GeoHashServer) GHMembers(_ context.Context, request *pb.GHMembersRequest) (*pb.GHMembersResponse, error) {
	res, err := service.GHMembers(request.Key)
	return &pb.GHMembersResponse{Points: res}, err
}
