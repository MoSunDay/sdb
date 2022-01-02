package fsm

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("GHCreate", func(logEntry *pb.LogEntry) interface{} {
		request := pb.GHCreateRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.GHCreate(request.Key, request.Precision)
		return &ApplyResponse{Response: &pb.GHCreateResponse{Success: res}, Err: err}
	})

	RegisterHandler("GHDel", func(logEntry *pb.LogEntry) interface{} {
		request := pb.GHDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.GHDel(request.Key)
		return &ApplyResponse{Response: &pb.GHDelResponse{Success: res}, Err: err}
	})

	RegisterHandler("GHAdd", func(logEntry *pb.LogEntry) interface{} {
		request := pb.GHAddRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.GHAdd(request.Key, request.Points)
		return &ApplyResponse{Response: &pb.GHAddResponse{Success: res}, Err: err}
	})

	RegisterHandler("GHRem", func(logEntry *pb.LogEntry) interface{} {
		request := pb.GHRemRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.GHRem(request.Key, request.Ids)
		return &ApplyResponse{Response: &pb.GHRemResponse{Success: res}, Err: err}
	})
}
