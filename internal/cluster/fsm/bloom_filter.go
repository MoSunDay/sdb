package fsm

import (
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("BFCreate", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.BFCreateRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.BFCreate(request.Key, request.N, request.P, batch)
	})

	RegisterHandler("BFDel", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.BFDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.BFDel(request.Key, batch)
	})

	RegisterHandler("BFAdd", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.BFAddRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.BFAdd(request.Key, request.Values, batch)
	})
}
