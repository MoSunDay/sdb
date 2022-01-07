package fsm

import (
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("Set", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.SetRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.Set(request.Key, request.Value, batch)
	})

	RegisterHandler("MSet", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.MSetRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.MSet(request.Keys, request.Values, batch)
	})

	RegisterHandler("SetNX", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.SetNXRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.SetNX(request.Key, request.Value, batch)
	})

	RegisterHandler("Del", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.DelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.Del(request.Key, batch)
	})

	RegisterHandler("Incr", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.IncrRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.Incr(request.Key, request.Delta, batch)
	})
}
