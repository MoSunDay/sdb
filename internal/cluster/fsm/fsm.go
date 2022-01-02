package fsm

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/raft"
	"github.com/yemingfeng/sdb/internal/pb"
	"google.golang.org/protobuf/proto"
	"io"
	log2 "log"
	"strconv"
	"sync"
)

var handlers = make(map[string]func(*pb.LogEntry) interface{}, 0)

func RegisterHandler(methodName string, handler func(*pb.LogEntry) interface{}) {
	handlers[methodName] = handler
}

type ApplyResponse struct {
	Response proto.Message
	Err      error
}

type FSM struct {
	sync.Mutex
	committed []string
}

func NewFSM() *FSM {
	return &FSM{
		committed: make([]string, 0),
	}
}

func (fsm *FSM) Apply(log *raft.Log) interface{} {
	logEntry := &pb.LogEntry{}
	err := proto.Unmarshal(log.Data, logEntry)
	if err != nil {
		return err
	}

	fsm.Lock()
	fsm.committed = append(fsm.committed, strconv.Itoa(int(log.Term))+":"+strconv.Itoa(int(log.Index)))
	defer fsm.Unlock()

	log2.Printf("apply: logEntry: [%+v], log: [%+v]", logEntry, log.Index)
	handler := handlers[logEntry.MethodName]
	if handler == nil {
		return &ApplyResponse{Err: errors.New("not support method name")}
	}
	return handler(logEntry)
}

type FSMSnapshot struct {
	committed []string
}

func (fsmSnapshot *FSMSnapshot) Persist(sink raft.SnapshotSink) error {
	bs, err := json.Marshal(fsmSnapshot.committed)
	if err != nil {
		return err
	}
	_, err = sink.Write(bs)
	return err
}

func (fsmSnapshot *FSMSnapshot) Release() {

}

func (fsm *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &FSMSnapshot{committed: fsm.committed}, nil
}

func (fsm *FSM) Restore(closer io.ReadCloser) error {
	bs, err := io.ReadAll(closer)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, &fsm.committed)
}
