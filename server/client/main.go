package main

import (
	"context"
	trippb "coolcar/proto/gen/go"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.SetFlags(log.Lshortfile)
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials())) // 建立不安全的TCP连接
	if err != nil {
		log.Fatalf("cannot connect server: %V", err)
	}

	tsClient := trippb.NewTripServiceClient(conn)
	r, err := tsClient.GetTrip(context.Background(), &trippb.GetTripRequest{Id: "456"})
	if err != nil {
		log.Fatalf("cannot call gettrip: %V", err)
	}

	fmt.Println(r)
}
