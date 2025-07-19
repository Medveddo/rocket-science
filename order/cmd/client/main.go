package main

import (
	"context"
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	orderV1 "github.com/Medveddo/rocket-science/shared/pkg/openapi/order/v1"
)

const (
	serverURL = "http://localhost:8080"
)

func main() {
	ctx := context.Background()

	// Инициализация Ogen-клиента
	client, err := orderV1.NewClient(serverURL)
	if err != nil {
		log.Fatalf("❌ Ошибка при создании клиента: %v", err)
	}

	response, err := client.CreateOrder(ctx, &orderV1.CreateOrderRequest{
		UserUUID: uuid.MustParse(gofakeit.UUID()),
		PartUuids: []string{
			"111e4567-e89b-12d3-a456-426614174001",
			"222e4567-e89b-12d3-a456-426614174002",
		},
	})
	if err != nil {
		log.Printf("❌ Ошибка при создании заказа: %v\n", err)
		return
	}
	log.Println(response)
	orderOK, ok := response.(*orderV1.CreateOrderResponse)
	if !ok {
		log.Printf("❌ Неожиданный тип ответа: %T\n", response)
		return
	}
	if orderOK.GetOrderUUID().String() == "" {
		log.Printf("❌ Путой UUID Заказа: %T\n", response)
		return
	}
	log.Println(orderOK.TotalPrice, "TotalPrice")
	if orderOK.TotalPrice == 0 {
		log.Printf("❌ Сумма заказа равна нулю: %v\n", orderOK.TotalPrice)
		return
	}

	pay_response, err := client.PayOrder(ctx, &orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PayOrderRequestPaymentMethodSBP,
	}, orderV1.PayOrderParams{
		OrderUUID: orderOK.OrderUUID.String(),
	})
	if err != nil {
		log.Printf("❌ Ошибка при оплате заказа: %v\n", err)
		return
	}
	payOK, ok := pay_response.(*orderV1.PayOrderResponse)
	if !ok {
		log.Printf("❌ Неожиданный тип ответа: %T\n", response)
		return
	}
	log.Println(payOK.TransactionUUID)

	getOrderResponse, err := client.GetOrder(ctx, orderV1.GetOrderParams{
		OrderUUID: orderOK.OrderUUID.String(),
	})
	if err != nil {
		log.Printf("❌ Ошибка при получкении заказа: %v\n", err)
		return
	}

	getOrderOK, ok := getOrderResponse.(*orderV1.OrderDto)
	if !ok {
		log.Printf("❌ Неожиданный тип ответа: %T\n", response)
		return
	}
	log.Println(getOrderOK.Status)
}
