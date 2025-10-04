package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// https://medium.com/@victor.neuret/mocking-the-official-mongo-golang-driver-5aad5b226a78
// https://github.com/gofr-dev/gofr/blob/development/pkg/gofr/datasource/mongo/mongo_test.go

func Test_AutoClosed_In_CreateAuction(t *testing.T) {

	os.Setenv("AUCTION_INTERVAL", "10ms")

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		auctionStatusMap := NewAuctionStatusMap()
		auctionRepository := NewAuctionRepository(mt.DB, auctionStatusMap)

		auction := &auction_entity.Auction{
			Id:          "1",
			ProductName: "Product 1",
			Category:    "Category 1",
			Description: "Description 1",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		err := auctionRepository.CreateAuction(context.Background(), auction)

		time.Sleep(50 * time.Millisecond)

		assert.Nil(t, err)

		status, okauction := auctionStatusMap.GetAuctionStatus(auction.Id)
		assert.True(t, okauction)
		assert.Equal(t, auction_entity.Completed, status)

	})

}
