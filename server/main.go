package main

import (
	"context"
	trippb "coolcar/proto/gen/go"
	trip "coolcar/proto/tripservice"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.SetFlags(log.Lshortfile)
	go startGRPCGateway()

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %V", err)
	}

	s := grpc.NewServer()
	trippb.RegisterTripServiceServer(s, &trip.Service{})
	log.Fatal(s.Serve(lis))
}

func startGRPCGateway() {
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel() // 函数返回后连接就会被cancel，下面建立的连接就会断开

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{
			EnumsAsInts: true, // 枚举类型输出整数
			OrigName:    true, // 返回原始名称，false会返回驼峰式命名
		},
	)) // 分发器，代理转发到多个微服务

	err := trippb.RegisterTripServiceHandlerFromEndpoint(
		c,                // context连接
		mux,              // 分发器，连接注册在mux
		"localhost:8081", // 连接端口，微服务
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, // tcp连接
	)

	if err != nil {
		log.Fatalf("cannot start grpc gateway: %v", err)
	}

	err = http.ListenAndServe(":8080", mux) // grpc gateway服务，mux式handler
	if err != nil {
		log.Fatalf("cannot listen and server: %v", err)
	}
}
