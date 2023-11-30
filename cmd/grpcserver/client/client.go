package main

import (
	"context"
	"go-users-service/cmd/grpcserver/usersproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("localhost:8081", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := usersproto.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user, err := client.SearchUser(ctx, &usersproto.SearchRequest{
		User: "6f35b5b0-1250-4b50-b58a-091d110cd00a",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user)
}
