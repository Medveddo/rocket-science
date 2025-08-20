//go:build integration

package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	Describe("GetPart", func() {
		var partUUID string

		BeforeEach(func() {
			// Вставляем тестовую деталь
			var err error
			partUUID, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовой детали в MongoDB")
		})

		It("должен успешно возвращать деталь по UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: partUUID,
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().Uuid).To(Equal(partUUID))
			Expect(resp.GetPart().GetName()).To(Equal("Test Engine Part"))
			Expect(resp.GetPart().GetDescription()).To(Equal("A test engine part for e2e testing"))
			Expect(resp.GetPart().GetPrice()).To(Equal(150000.00))
			Expect(resp.GetPart().GetStockQuantity()).To(Equal(int64(10)))
			Expect(resp.GetPart().GetCategory()).To(Equal(inventoryV1.Category_CATEGORY_ENGINE))
			Expect(resp.GetPart().GetCreatedAt()).ToNot(BeNil())
			Expect(resp.GetPart().GetUpdatedAt()).ToNot(BeNil())
		})
	})

	Describe("ListParts", func() {
		BeforeEach(func() {
			// Вставляем несколько тестовых деталей
			_, err := env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку первой тестовой детали")

			_, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку второй тестовой детали")
		})

		It("должен успешно возвращать список всех деталей", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeNil())
			Expect(len(resp.GetParts())).To(BeNumerically(">=", 2))

			// Проверяем, что все детали имеют необходимые поля
			for _, part := range resp.GetParts() {
				Expect(part.GetUuid()).ToNot(BeEmpty())
				Expect(part.GetName()).ToNot(BeEmpty())
				Expect(part.GetDescription()).ToNot(BeEmpty())
				Expect(part.GetPrice()).To(BeNumerically(">", 0))
				Expect(part.GetStockQuantity()).To(BeNumerically(">=", 0))
				Expect(part.GetCreatedAt()).ToNot(BeNil())
				Expect(part.GetUpdatedAt()).ToNot(BeNil())
			}
		})
	})
})
