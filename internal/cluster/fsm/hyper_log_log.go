package fsm

import (
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("HLLCreate", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.HLLCreateRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.HLLCreate(request.Key, batch)
	})

	RegisterHandler("HLLDel", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.HLLDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.HLLDel(request.Key, batch)
	})

	RegisterHandler("HLLAdd", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.HLLAddRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.HLLAdd(request.Key, request.Values, batch)
	})
}
