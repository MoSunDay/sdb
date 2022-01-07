package fsm

import (
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("ZPush", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.ZPushRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.ZPush(request.Key, request.Tuples, batch)
	})

	RegisterHandler("ZPop", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.ZPopRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.ZPop(request.Key, request.Values, batch)
	})

	RegisterHandler("ZDel", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.ZDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.ZDel(request.Key, batch)
	})
}
