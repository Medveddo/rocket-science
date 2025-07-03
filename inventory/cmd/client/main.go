package main

import (
	"context"
	"log"

	// "github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
)

const serverAddress = "localhost:50051"

// func getPart(ctx context.Context, client inventoryV1.InventoryServiceClient, uuid string) (*inventoryV1.Part, error) {
// 	resp, err := client.GetPart(ctx, &inventoryV1.GetPartRequest{Uuid: uuid})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp.Part, nil
// }

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
	client := inventoryV1.NewInventoryServiceClient(conn)
	// part_uuid := "111e4567-e89b-12d3-a456-426614174001"

	resp, err := client.ListParts(ctx, &inventoryV1.ListPartsRequest{Filter: &inventoryV1.PartsFilter{
		// Uuids: []string{part_uuid},
		// Names: []string{"Quantum Shield Generator"},
		ManufacturerCountries: []string{"USA"},
	}})
	if err != nil {
		log.Printf("listPart error: %v\n", err)
		return
	}
	for _, part := range resp.Parts {
		log.Println("part: ", part)
	}
	// part, err := getPart(ctx, client, part_uuid)
	// if err != nil {
	// 	log.Printf("getPart error: %v\n", err)
	// 	return
	// }

	// Выводим информацию о полученном наблюдении
	// log.Printf("Got part: UUID=%s", part_uuid)
	// log.Printf("%v\n", part)
}
