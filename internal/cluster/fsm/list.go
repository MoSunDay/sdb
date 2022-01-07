package fsm

import (
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("LRPush", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.LRPushRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.LRPush(request.Key, request.Values, batch)
	})

	RegisterHandler("LLPush", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.LLPushRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.LLPush(request.Key, request.Values, batch)
	})

	RegisterHandler("LPop", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.LPopRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.LPop(request.Key, request.Values, batch)
	})

	RegisterHandler("LDel", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.LDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.LDel(request.Key, batch)
	})
}
