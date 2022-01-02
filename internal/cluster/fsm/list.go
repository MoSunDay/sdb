package fsm

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("LRPush", func(logEntry *pb.LogEntry) interface{} {
		request := pb.LRPushRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.LRPush(request.Key, request.Values)
		return &ApplyResponse{Response: &pb.LRPushResponse{Success: res}, Err: err}
	})

	RegisterHandler("LLPush", func(logEntry *pb.LogEntry) interface{} {
		request := pb.LLPushRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.LLPush(request.Key, request.Values)
		return &ApplyResponse{Response: &pb.LLPushResponse{Success: res}, Err: err}
	})

	RegisterHandler("LPop", func(logEntry *pb.LogEntry) interface{} {
		request := pb.LPopRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.LPop(request.Key, request.Values)
		return &ApplyResponse{Response: &pb.LPopResponse{Success: res}, Err: err}
	})

	RegisterHandler("LDel", func(logEntry *pb.LogEntry) interface{} {
		request := pb.LDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.LDel(request.Key)
		return &ApplyResponse{Response: &pb.LDelResponse{Success: res}, Err: err}
	})
}
