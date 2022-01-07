package fsm

import (
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("BSCreate", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.BSCreateRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.BSCreate(request.Key, request.Size, batch)
	})

	RegisterHandler("BSDel", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.BSDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.BSDel(request.Key, batch)
	})

	RegisterHandler("BSSetRange", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.BSSetRangeRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.BSSetRange(request.Key, request.Start, request.End, request.Value, batch)
	})

	RegisterHandler("BSMSet", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.BSMSetRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.BSMSet(request.Key, request.Bits, request.Value, batch)
	})
}
