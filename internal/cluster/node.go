package cluster

import (
	"errors"
	"fmt"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/yemingfeng/sdb/internal/cluster/fsm"
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/pb"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var node *raft.Raft

func Start() {
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(conf.Conf.Cluster.NodeId)
	raftConfig.SnapshotThreshold = 1
	raftConfig.SnapshotInterval = 1 * time.Second
	log.Printf("raft config: %+v", raftConfig)

	path := filepath.Join(conf.Conf.Cluster.Path, string(raftConfig.LocalID))

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatalln("mkdir error", err)
	}

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(path, "raft-log.bolt"))
	if err != nil {
		log.Fatalf("new log store error, %+v", err)
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(path, "raft-stable.bolt"))
	if err != nil {
		log.Fatalf("new stable store error, %+v", err)
	}
	snapshotStore, err := raft.NewFileSnapshotStore(path, 1, os.Stderr)
	if err != nil {
		log.Fatalln("new snapshot store error", err)
	}

	address, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+strconv.Itoa(conf.Conf.Cluster.Port))
	if err != nil {
		log.Fatalln("resolve tcp addr error", err)
	}
	transport, err := raft.NewTCPTransport(address.String(), address, 3, 10*time.Second, os.Stderr)
	if err != nil {
		log.Fatalln("new tcp transport addr error", err)
	}

	node, err = raft.NewRaft(raftConfig, fsm.NewFSM(), logStore, stableStore, snapshotStore, transport)
	if err != nil {
		log.Fatalln("new raft error", err)
	}

	if len(conf.Conf.Cluster.Master) == 0 {
		node.BootstrapCluster(raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      raftConfig.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		})
	} else {
		url := fmt.Sprintf("http://%s/join?nodeId=%s&joinAddress=%s", conf.Conf.Cluster.Master, raftConfig.LocalID, transport.LocalAddr())
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("join cluster error: %+v", err)
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("read body error: %+v", err)
		}
		log.Printf("join response: %s", body)
		if string(body) != "ok" {
			log.Fatalf("join cluster error: %s", body)
		}
	}
}

func Apply(methodName string, request proto.Message) (proto.Message, error) {
	if node.State() != raft.Leader {
		return nil, errors.New("only can apply from leader")
	}
	requestBytes, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}
	entry := &pb.LogEntry{MethodName: methodName, RequestBytes: requestBytes}
	data, err := proto.Marshal(entry)
	if err != nil {
		return nil, err
	}

	future := node.Apply(data, time.Duration(conf.Conf.Cluster.Timeout)*time.Millisecond)
	if err := future.Error(); err != nil {
		log.Fatalf("apply error: %+v", err)
		return nil, err
	}
	applyResponse := future.Response().(*fsm.ApplyResponse)
	return applyResponse.Response, applyResponse.Err
}

func Join(nodeID, addr string) error {
	log.Printf("received join request for remote node %s at %s", nodeID, addr)

	configFuture := node.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		log.Println("failed to get raft configuration", err)
		return err
	}

	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeID) {
				log.Printf("node %s at %s already member of cluster, ignoring join request", nodeID, addr)
				return nil
			}

			future := node.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", nodeID, addr, err)
			}
		}
	}
	log.Println("Coming to add voter")
	f := node.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		log.Println(f.Error())
		return f.Error()
	}
	log.Printf("node %s at %s joined successfully", nodeID, addr)
	return nil
}
