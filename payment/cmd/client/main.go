package main

import (
	"context"
	"log"

	fake "github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	paymentV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/payment/v1"
)

const serverAddress = "localhost:50052"

func main() {
	ctx := context.Background()

	conn, err := grpc.NewClient(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	// Создаем gRPC клиент
	client := paymentV1.NewPaymentServiceClient(conn)

	response, err := client.PayOrder(ctx, &paymentV1.PayOrderRequest{
		UserUuid:      fake.UUID(),
		OrderUuid:     fake.UUID(),
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
	})
	if err != nil {
		log.Println("error while PayOrder:", err)
	}
	log.Println(response)
}
