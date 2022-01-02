package fsm

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("HLLCreate", func(logEntry *pb.LogEntry) interface{} {
		request := pb.HLLCreateRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.HLLCreate(request.Key)
		return &ApplyResponse{Response: &pb.HLLCreateResponse{Success: res}, Err: err}
	})

	RegisterHandler("HLLDel", func(logEntry *pb.LogEntry) interface{} {
		request := pb.HLLDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.HLLDel(request.Key)
		return &ApplyResponse{Response: &pb.HLLDelResponse{Success: res}, Err: err}
	})

	RegisterHandler("HLLAdd", func(logEntry *pb.LogEntry) interface{} {
		request := pb.HLLAddRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.HLLAdd(request.Key, request.Values)
		return &ApplyResponse{Response: &pb.HLLAddResponse{Success: res}, Err: err}
	})
}
