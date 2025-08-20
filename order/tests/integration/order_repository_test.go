//go:build integration

package integration

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Medveddo/rocket-science/order/internal/model"
)

var _ = Describe("OrderRepository", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = suiteCtx

		// Clear the orders table before each test
		err := env.ClearOrdersTable(ctx)
		Expect(err).ToNot(HaveOccurred(), "expected to successfully clear orders table")
	})

	Describe("Get", func() {
		Context("when order does not exist", func() {
			It("should return ErrOrderNotFound", func() {
				nonExistentUUID := uuid.New().String()
				_, err := env.Repo.Get(ctx, nonExistentUUID)

				Expect(err).To(HaveOccurred(), "expected error when getting non-existent order")
				Expect(err).To(Equal(model.ErrOrderNotFound))
			})
		})
	})

	Describe("Update and Get", func() {
		It("should insert order and return consistent data", func() {
			// Create test order data
			orderUUID := uuid.New()
			userUUID := uuid.New()
			partUUIDs := []string{uuid.New().String(), uuid.New().String()}
			totalPrice := 150000.00
			status := model.OrderStatusPENDINGPAYMENT

			order := model.Order{
				OrderUUID:       orderUUID,
				UserUUID:        userUUID,
				PartUUIDs:       partUUIDs,
				TotalPrice:      totalPrice,
				TransactionUUID: nil,
				PaymentMethod:   nil,
				Status:          status,
			}

			// Insert order via repository
			err := env.Repo.Update(ctx, order)
			Expect(err).ToNot(HaveOccurred(), "expected no error when inserting order")

			// Retrieve order via repository
			retrievedOrder, err := env.Repo.Get(ctx, orderUUID.String())
			Expect(err).ToNot(HaveOccurred(), "expected no error when getting order")

			// Assert data consistency
			Expect(retrievedOrder.OrderUUID).To(Equal(orderUUID))
			Expect(retrievedOrder.UserUUID).To(Equal(userUUID))
			Expect(retrievedOrder.PartUUIDs).To(Equal(partUUIDs))
			Expect(retrievedOrder.TotalPrice).To(Equal(totalPrice))
			Expect(retrievedOrder.Status).To(Equal(status))
			Expect(retrievedOrder.TransactionUUID).To(BeNil())
			Expect(retrievedOrder.PaymentMethod).To(BeNil())
		})

		It("should update existing order and return consistent data", func() {
			// First insert an order
			orderUUID := uuid.New()
			userUUID := uuid.New()
			partUUIDs := []string{uuid.New().String()}
			totalPrice := 50000.00
			status := model.OrderStatusPENDINGPAYMENT

			order := model.Order{
				OrderUUID:       orderUUID,
				UserUUID:        userUUID,
				PartUUIDs:       partUUIDs,
				TotalPrice:      totalPrice,
				TransactionUUID: nil,
				PaymentMethod:   nil,
				Status:          status,
			}

			err := env.Repo.Update(ctx, order)
			Expect(err).ToNot(HaveOccurred(), "expected no error when inserting order")

			// Update the order
			updatedOrder := order
			updatedOrder.Status = model.OrderStatusPAID
			updatedOrder.TotalPrice = 75000.00
			transactionUUID := uuid.New()
			paymentMethod := model.PayOrderRequestPaymentMethodCARD
			updatedOrder.TransactionUUID = &transactionUUID
			updatedOrder.PaymentMethod = &paymentMethod

			err = env.Repo.Update(ctx, updatedOrder)
			Expect(err).ToNot(HaveOccurred(), "expected no error when updating order")

			// Retrieve updated order and verify consistency
			retrievedOrder, err := env.Repo.Get(ctx, orderUUID.String())
			Expect(err).ToNot(HaveOccurred(), "expected no error when getting updated order")

			// Assert updated data consistency
			Expect(retrievedOrder.OrderUUID).To(Equal(orderUUID))
			Expect(retrievedOrder.UserUUID).To(Equal(userUUID))
			Expect(retrievedOrder.PartUUIDs).To(Equal(partUUIDs))
			Expect(retrievedOrder.TotalPrice).To(Equal(75000.00))
			Expect(retrievedOrder.Status).To(Equal(model.OrderStatusPAID))
			Expect(retrievedOrder.TransactionUUID).ToNot(BeNil())
			Expect(*retrievedOrder.TransactionUUID).To(Equal(transactionUUID))
			Expect(retrievedOrder.PaymentMethod).ToNot(BeNil())
			Expect(*retrievedOrder.PaymentMethod).To(Equal(paymentMethod))
		})
	})
})
