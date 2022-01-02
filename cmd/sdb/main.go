package main

import (
	"github.com/yemingfeng/sdb/internal/cluster"
	"github.com/yemingfeng/sdb/internal/server"
)

func main() {
	cluster.Start()

	httpServer := server.NewHttpServer()
	go func() {
		httpServer.Start()
	}()

	sdbGrpcServer := server.NewSDBGrpcServer()
	sdbGrpcServer.Start()
}
