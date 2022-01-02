package fsm

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("MPush", func(logEntry *pb.LogEntry) interface{} {
		request := pb.MPushRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.MPush(request.Key, request.Pairs)
		return &ApplyResponse{Response: &pb.MPushResponse{Success: res}, Err: err}
	})

	RegisterHandler("MPop", func(logEntry *pb.LogEntry) interface{} {
		request := pb.MPopRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.MPop(request.Key, request.Keys)
		return &ApplyResponse{Response: &pb.MPopResponse{Success: res}, Err: err}
	})

	RegisterHandler("MDel", func(logEntry *pb.LogEntry) interface{} {
		request := pb.MDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.MDel(request.Key)
		return &ApplyResponse{Response: &pb.MDelResponse{Success: res}, Err: err}
	})
}
