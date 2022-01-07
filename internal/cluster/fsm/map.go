package fsm

import (
	"github.com/yemingfeng/sdb/internal/engine"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/service"
	"google.golang.org/protobuf/proto"
)

func init() {
	RegisterHandler("MPush", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.MPushRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.MPush(request.Key, request.Pairs, batch)
	})

	RegisterHandler("MPop", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.MPopRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.MPop(request.Key, request.Keys, batch)
	})

	RegisterHandler("MDel", func(logEntry *pb.LogEntry, batch engine.Batch) error {
		request := pb.MDelRequest{}
		err := proto.Unmarshal(logEntry.RequestBytes, &request)
		if err != nil {
			return err
		}
		return service.MDel(request.Key, batch)
	})
}
