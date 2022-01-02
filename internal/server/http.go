package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yemingfeng/sdb/internal/cluster"
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/pb"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type HttpServer struct {
	mux *runtime.ServeMux
}

func NewHttpServer() *HttpServer {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterSDBHandlerFromEndpoint(context.Background(),
		mux, ":"+strconv.Itoa(conf.Conf.Server.GRPCPort), opts)
	if err != nil {
		log.Fatalf("failed to register: %+v", err)
	}
	return &HttpServer{mux: mux}
}

func (httpServer *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.RequestURI == "/metrics" {
		promhttp.Handler().ServeHTTP(writer, request)
	} else if strings.HasPrefix(request.RequestURI, "/join") {
		nodeId := request.URL.Query()["nodeId"][0]
		joinAddress := request.URL.Query()["joinAddress"][0]
		err := cluster.Join(nodeId, joinAddress)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			writer.WriteHeader(502)
			return
		}
		_, err = writer.Write([]byte("ok"))
		if err != nil {
			log.Println(err)
			writer.WriteHeader(502)
		}
	} else if strings.HasPrefix(request.RequestURI, "/v1") {
		httpServer.mux.ServeHTTP(writer, request)
	} else {
		writer.WriteHeader(502)
	}
}

func (httpServer *HttpServer) Start() {
	err := http.ListenAndServe(
		":"+strconv.Itoa(conf.Conf.Server.HttpPort), httpServer)
	if err != nil {
		log.Fatalf("failed to serve: %+v", err)
	}
}
