package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/server"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	lg, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create zap logger: %v", err)
	}

	c := context.Background()
	c, cancel := context.WithCancel(c)

	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{
			EnumsAsInts: true,
			OrigName:    true,
		},
	))

	serverConfig := []struct {
		name         string
		addr         string
		registerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	}{
		{
			name:         "auth",
			addr:         "localhost:8081",
			registerFunc: authpb.RegisterAuthServiceHandlerFromEndpoint,
		},
		{
			name:         "rental",
			addr:         "localhost:8082",
			registerFunc: rentalpb.RegisterTripServiceHandlerFromEndpoint,
		},
	}

	for _, s := range serverConfig {
		err := s.registerFunc(
			c,
			mux,
			s.addr,
			[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		)
		if err != nil {
			lg.Sugar().Fatalf("cannot start grpc gateway: %s: %v", s.name, err)
		}
	}

	// err := authpb.RegisterAuthServiceHandlerFromEndpoint(
	// 	c,
	// 	mux,
	// 	"localhost:8081",
	// 	[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	// )
	// if err != nil {
	// 	log.Fatalf("cannot start grpc gateway: %v", err)
	// }

	// err = rentalpb.RegisterTripServiceHandlerFromEndpoint(
	// 	c,
	// 	mux,
	// 	"localhost:8082",
	// 	[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	// )
	// if err != nil {
	// 	log.Fatalf("cannot start grpc gateway: %v", err)
	// }

	addr := ":8080"
	lg.Sugar().Infof("grpc gateway started at %s", addr)
	lg.Sugar().Fatal(http.ListenAndServe(addr, mux))
}
