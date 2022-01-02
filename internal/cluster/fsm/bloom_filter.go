package fsm

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("BFCreate", func(logEntry *pb.LogEntry) interface{} {
		request := pb.BFCreateRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.BFCreate(request.Key, request.N, request.P)
		return &ApplyResponse{Response: &pb.BFCreateResponse{Success: res}, Err: err}
	})

	RegisterHandler("BFDel", func(logEntry *pb.LogEntry) interface{} {
		request := pb.BFDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.BFDel(request.Key)
		return &ApplyResponse{Response: &pb.BFDelResponse{Success: res}, Err: err}
	})

	RegisterHandler("BFAdd", func(logEntry *pb.LogEntry) interface{} {
		request := pb.BFAddRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.BFAdd(request.Key, request.Values)
		return &ApplyResponse{Response: &pb.BFAddResponse{Success: res}, Err: err}
	})
}
