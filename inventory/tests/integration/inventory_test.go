//go:build integration

package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		conn            *grpc.ClientConn
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		var err error
		conn, err = grpc.DialContext(
			ctx,
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

		_ = conn.Close()
		cancel()
	})

	Describe("GetPart", func() {
		It("должен успешно возвращать деталь по UUID", func() {
			// Подготовка данных: создаём и сохраняем деталь в БД
			part, err := env.GetTestPart(ctx)
			Expect(err).ToNot(HaveOccurred())

			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: part.Uuid,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())

			got := resp.GetPart()
			Expect(got.GetUuid()).To(Equal(part.Uuid))
			Expect(got.GetName()).To(Equal(part.Name))
			Expect(got.GetDescription()).To(Equal(part.Description))
			Expect(got.GetPrice()).To(Equal(part.Price))
			Expect(got.GetStockQuantity()).To(Equal(part.StockQuantity))
			Expect(got.GetCategory().String()).ToNot(BeEmpty())
			Expect(got.GetDimensions()).ToNot(BeNil())
			Expect(got.GetManufacturer()).ToNot(BeNil())
			Expect(len(got.GetTags())).To(BeNumerically(">=", 0))
		})
	})

	Describe("ListParts", func() {
		It("должен возвращать список деталей (как минимум одну)", func() {
			// Подготовим данные
			_, err := env.GetTestPart(ctx)
			Expect(err).ToNot(HaveOccurred())

			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{}, // без фильтров
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeNil())
			Expect(len(resp.GetParts())).To(BeNumerically(">=", 1))
		})
	})
})
