package server

import (
	pb2 "github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"golang.org/x/net/context"
)

type GeoHashServer struct {
	pb2.UnimplementedSDBServer
}

func (server *GeoHashServer) GHCreate(_ context.Context, request *pb2.GHCreateRequest) (*pb2.GHCreateResponse, error) {
	res, err := service.GHCreate(request.Key, request.Precision)
	return &pb2.GHCreateResponse{Success: res}, err
}

func (server *GeoHashServer) GHDel(_ context.Context, request *pb2.GHDelRequest) (*pb2.GHDelResponse, error) {
	res, err := service.GHDel(request.Key)
	return &pb2.GHDelResponse{Success: res}, err
}

func (server *GeoHashServer) GHAdd(_ context.Context, request *pb2.GHAddRequest) (*pb2.GHAddResponse, error) {
	res, err := service.GHAdd(request.Key, request.Points)
	return &pb2.GHAddResponse{Success: res}, err
}

func (server *GeoHashServer) GHRem(_ context.Context, request *pb2.GHRemRequest) (*pb2.GHRemResponse, error) {
	res, err := service.GHRem(request.Key, request.Points)
	return &pb2.GHRemResponse{Success: res}, err
}

func (server *GeoHashServer) GHGetBoxes(_ context.Context, request *pb2.GHGetBoxesRequest) (*pb2.GHGetBoxesResponse, error) {
	res, err := service.GHGetBoxes(request.Key, request.Latitude, request.Longitude)
	return &pb2.GHGetBoxesResponse{Points: res}, err
}

func (server *GeoHashServer) GHGetNeighbors(_ context.Context, request *pb2.GHGetNeighborsRequest) (*pb2.GHGetNeighborsResponse, error) {
	res, err := service.GHGetNeighbors(request.Key, request.Latitude, request.Longitude)
	return &pb2.GHGetNeighborsResponse{Points: res}, err
}

func (server *GeoHashServer) GHCount(_ context.Context, request *pb2.GHCountRequest) (*pb2.GHCountResponse, error) {
	res, err := service.GHCount(request.Key)
	return &pb2.GHCountResponse{Count: res}, err
}
func (server *GeoHashServer) GHMembers(_ context.Context, request *pb2.GHMembersRequest) (*pb2.GHMembersResponse, error) {
	res, err := service.GHMembers(request.Key)
	return &pb2.GHMembersResponse{Points: res}, err
}
