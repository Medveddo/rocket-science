package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panicf("did not connect: %v", err)
		return
	}
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: "payment",
	})
	if err != nil {
		log.Panicf("failed to check health: %v", err)
		return
	}

	fmt.Printf("health status: %v\n", response.Status)
}
