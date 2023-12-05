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
	conn, err := grpc.Dial("localhost:3030", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := usersproto.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user, err := client.SearchUser(ctx, &usersproto.SearchRequest{
		User: "c249e993-9b2d-4b75-b0b6-f462dbf2666e",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user)
}
