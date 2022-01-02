package fsm

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("ZPush", func(logEntry *pb.LogEntry) interface{} {
		request := pb.ZPushRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.ZPush(request.Key, request.Tuples)
		return &ApplyResponse{Response: &pb.ZPushResponse{Success: res}, Err: err}
	})

	RegisterHandler("ZPop", func(logEntry *pb.LogEntry) interface{} {
		request := pb.ZPopRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.ZPop(request.Key, request.Values)
		return &ApplyResponse{Response: &pb.ZPopResponse{Success: res}, Err: err}
	})

	RegisterHandler("ZDel", func(logEntry *pb.LogEntry) interface{} {
		request := pb.ZDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.ZDel(request.Key)
		return &ApplyResponse{Response: &pb.ZDelResponse{Success: res}, Err: err}
	})
}
