package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/dao"
	"coolcar/auth/token"
	"coolcar/auth/wechat"
	"coolcar/shared/server"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}

	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://admin:123456@localhost:27017/?authSource=admin&readPreference=primary&ssl=false&directConnection=true"))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}

	pkFile, err := os.Open("private.key")
	if err != nil {
		logger.Fatal("cannot open private key", zap.Error(err))
	}

	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		logger.Fatal("cannot read private key", zap.Error(err))
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		logger.Fatal("cannot parse private key", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:   "auth",
		Addr:   ":8081",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			authpb.RegisterAuthServiceServer(s, &auth.Service{
				OpenIdResolver: &wechat.Service{
					AppID:     "wx7c22d4fc3c83c721",               // wxcf5e3e5314e34ace
					AppSecret: "d6d7f52eba82a0ff0aa53f262ef5f5dc", // AppID(小程序ID)wx7c22d4fc3c83c721
				},
				Mongo:          dao.NewMongo(mongoClient.Database("coolcar")),
				Logger:         *logger,
				TokenExpire:    2 * time.Hour,
				TokenGenerator: token.NewJWTTokenGen("coolcar/auth", privKey),
			})
		},
	}))
}

// func main() {
// 	// logger, err := zap.NewDevelopment()
// 	logger, err := newZapLogger()

// 	if err != nil {
// 		log.Fatalf("cannot create logger: %v", err)
// 	}

// 	lis, err := net.Listen("tcp", ":8081")
// 	if err != nil {
// 		logger.Fatal("cannot listen", zap.Error(err))
// 	}

// 	c := context.Background()
// 	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://admin:123456@localhost:27017/?authSource=admin&readPreference=primary&ssl=false&directConnection=true"))
// 	if err != nil {
// 		logger.Fatal("cannot connect mongodb", zap.Error(err))
// 	}

// 	pkFile, err := os.Open("private.key")
// 	if err != nil {
// 		logger.Fatal("cannot open private key", zap.Error(err))
// 	}

// 	pkBytes, err := ioutil.ReadAll(pkFile)
// 	if err != nil {
// 		logger.Fatal("cannot read private key", zap.Error(err))
// 	}

// 	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
// 	if err != nil {
// 		logger.Fatal("cannot parse private key", zap.Error(err))
// 	}

// 	s := grpc.NewServer()
// 	authpb.RegisterAuthServiceServer(s, &auth.Service{
// 		OpenIdResolver: &wechat.Service{
// 			AppID:     "wx7c22d4fc3c83c721",               // wxcf5e3e5314e34ace
// 			AppSecret: "d6d7f52eba82a0ff0aa53f262ef5f5dc", // AppID(小程序ID)wx7c22d4fc3c83c721
// 		},
// 		Mongo:          dao.NewMongo(mongoClient.Database("coolcar")),
// 		Logger:         *logger,
// 		TokenExpire:    2 * time.Hour,
// 		TokenGenerator: token.NewJWTTokenGen("coolcar/auth", privKey),
// 	})

// 	err = s.Serve(lis)
// 	logger.Fatal("cannot server", zap.Error(err))
// }

// // 自定义zapLogger
// func newZapLogger() (*zap.Logger, error) {
// 	cfg := zap.NewDevelopmentConfig()
// 	cfg.EncoderConfig.TimeKey = ""
// 	return cfg.Build()
// }
