package fsm

import (
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("GHCreate", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.GHCreateRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.GHCreate(request.Key, request.Precision, batch)
	})

	RegisterHandler("GHDel", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.GHDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.GHDel(request.Key, batch)
	})

	RegisterHandler("GHAdd", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.GHAddRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.GHAdd(request.Key, request.Points, batch)
	})

	RegisterHandler("GHRem", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.GHRemRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.GHRem(request.Key, request.Ids, batch)
	})
}
