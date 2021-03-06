package main

import (
	"net"

	"github.com/grpcbrick/queues/database"
	"github.com/grpcbrick/queues/provider"
	"github.com/grpcbrick/queues/standard"
	"github.com/yinxulai/goutils/config"
	"github.com/yinxulai/goutils/grpc/interceptor"
	"google.golang.org/grpc"
)

// 队列不是一个高读写并发的东西

func init() {
	config.SetStandard("rpc_port", "", true, "RPC 使用的 MYSQL 数据库配置")
	config.SetStandard("mysql_url", ":3000", true, "RPC 服务监听的端口")
}

func main() {
	database.Init()

	lis, err := net.Listen("tcp", config.MustGet("port"))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(interceptor.NewCalllogs()...)
	standard.RegisterQueuesServer(grpcServer, provider.NewService())

	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}
