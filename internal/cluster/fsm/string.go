package fsm

import (
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("Set", func(logEntry *pb.LogEntry) interface{} {
		request := pb.SetRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.Set(request.Key, request.Value)
		return &ApplyResponse{Response: &pb.SetResponse{Success: res}, Err: err}
	})

	RegisterHandler("MSet", func(logEntry *pb.LogEntry) interface{} {
		request := pb.MSetRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.MSet(request.Keys, request.Values)
		return &ApplyResponse{Response: &pb.MSetResponse{Success: res}, Err: err}
	})

	RegisterHandler("SetNX", func(logEntry *pb.LogEntry) interface{} {
		request := pb.SetNXRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.SetNX(request.Key, request.Value)
		return &ApplyResponse{Response: &pb.SetNXResponse{Success: res}, Err: err}
	})

	RegisterHandler("SetGet", func(logEntry *pb.LogEntry) interface{} {
		request := pb.SetGetRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, old, err := service.SetGet(request.Key, request.Value)
		return &ApplyResponse{Response: &pb.SetGetResponse{Success: res, OldValue: old}, Err: err}
	})

	RegisterHandler("Del", func(logEntry *pb.LogEntry) interface{} {
		request := pb.DelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.Del(request.Key)
		return &ApplyResponse{Response: &pb.DelResponse{Success: res}, Err: err}
	})

	RegisterHandler("Incr", func(logEntry *pb.LogEntry) interface{} {
		request := pb.IncrRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return &ApplyResponse{Err: err}
		}
		res, err := service.Incr(request.Key, request.Delta)
		return &ApplyResponse{Response: &pb.IncrResponse{Success: res}, Err: err}
	})
}
