//go:build integration

package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)

var _ = Describe("Inventory", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), testsTimeout)
		conn, err := grpc.NewClient(env.App.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное очищение коллекции частей")
		cancel()
	})

	Describe("GetPart", func() {
		It("should return a part", func() {
			ctx := context.Background()
			partUUID, err := env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешное добавление части")
			Expect(partUUID).To(MatchRegexp(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`))

			part, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{Uuid: partUUID})
			Expect(err).ToNot(HaveOccurred(), "ожидали успешное получение части")
			Expect(part.Part.Uuid).To(Equal(partUUID), "ожидали получить ту же часть, что и добавленную")
		})
	})

	Describe("GetListParts", func() {
		It("should return a list of parts", func() {
			ctx := context.Background()
			partUUIDOne, err := env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешное добавление части")
			Expect(partUUIDOne).To(MatchRegexp(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`))

			partUUIDTwo, err := env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешное добавление части")
			Expect(partUUIDTwo).To(MatchRegexp(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`))

			listParts, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{})
			Expect(err).ToNot(HaveOccurred(), "ожидали успешное получение списка частей")
			Expect(listParts.Parts).To(HaveLen(2), "ожидали получить список из 2 частей")
			Expect(listParts.Parts[0].Uuid).To(Equal(partUUIDOne), "ожидали получить первую часть")
			Expect(listParts.Parts[1].Uuid).To(Equal(partUUIDTwo), "ожидали получить вторую часть")
		})
	})
})
