package fsm

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("SPush", func(logEntry *pb.LogEntry) interface{} {
		request := pb.SPushRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.SPush(request.Key, request.Values)
		return &ApplyResponse{Response: &pb.SPushResponse{Success: res}, Err: err}
	})

	RegisterHandler("SPop", func(logEntry *pb.LogEntry) interface{} {
		request := pb.SPopRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.SPop(request.Key, request.Values)
		return &ApplyResponse{Response: &pb.SPopResponse{Success: res}, Err: err}
	})

	RegisterHandler("SDel", func(logEntry *pb.LogEntry) interface{} {
		request := pb.SDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.SDel(request.Key)
		return &ApplyResponse{Response: &pb.SDelResponse{Success: res}, Err: err}
	})
}
