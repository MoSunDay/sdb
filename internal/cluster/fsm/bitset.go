package fsm

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("BSCreate", func(logEntry *pb.LogEntry) interface{} {
		request := pb.BSCreateRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.BSCreate(request.Key, request.Size)
		return &ApplyResponse{Response: &pb.BSCreateResponse{Success: res}, Err: err}
	})

	RegisterHandler("BSDel", func(logEntry *pb.LogEntry) interface{} {
		request := pb.BSDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.BSDel(request.Key)
		return &ApplyResponse{Response: &pb.BSDelResponse{Success: res}, Err: err}
	})

	RegisterHandler("BSSetRange", func(logEntry *pb.LogEntry) interface{} {
		request := pb.BSSetRangeRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.BSSetRange(request.Key, request.Start, request.End, request.Value)
		return &ApplyResponse{Response: &pb.BSSetRangeResponse{Success: res}, Err: err}
	})

	RegisterHandler("BSMSet", func(logEntry *pb.LogEntry) interface{} {
		request := pb.BSMSetRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.BSMSet(request.Key, request.Bits, request.Value)
		return &ApplyResponse{Response: &pb.BSMSetResponse{Success: res}, Err: err}
	})
}
