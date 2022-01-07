package fsm

import (
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("SPush", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.SPushRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.SPush(request.Key, request.Values, batch)
	})

	RegisterHandler("SPop", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.SPopRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.SPop(request.Key, request.Values, batch)
	})

	RegisterHandler("SDel", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.SDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.SDel(request.Key, batch)
	})
}
