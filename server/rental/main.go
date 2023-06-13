package main

import (
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip"
	"coolcar/shared/server"
	"log"

	"google.golang.org/grpc"
)

func main() {
	// logger, err := zap.NewDevelopment()
	logger, err := server.NewZapLogger()

	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	// 或者err = server.RunGRPCServer() .Sugar().Fatal是语法糖
	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Addr:              ":8082",
		AuthPublicKeyFile: "../shared/public.key",
		Logger:            logger,
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				Logger: logger,
			})
		},
	}))

	// lis, err := net.Listen("tcp", ":8082")
	// if err != nil {
	// 	logger.Fatal("cannot listen", zap.Error(err))
	// }

	// in, err := auth.Interceptor("../shared/public.key")
	// if err != nil {
	// 	logger.Fatal("cannot create auth interceptor", zap.Error(err))
	// }

	// // 定义interceptor
	// s := grpc.NewServer(grpc.UnaryInterceptor(in))
	// rentalpb.RegisterTripServiceServer(s, &trip.Service{
	// 	Logger: logger,
	// })

	// err = s.Serve(lis)
	// logger.Fatal("cannot server", zap.Error(err))
}

// 自定义zapLogger
// func newZapLogger() (*zap.Logger, error) {
// 	cfg := zap.NewDevelopmentConfig()
// 	cfg.EncoderConfig.TimeKey = ""
// 	return cfg.Build()
// }
