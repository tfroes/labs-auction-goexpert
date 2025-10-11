package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// https://medium.com/@victor.neuret/mocking-the-official-mongo-golang-driver-5aad5b226a78
// https://github.com/gofr-dev/gofr/blob/development/pkg/gofr/datasource/mongo/mongo_test.go

func Test_AutoClosed_In_CreateAuction(t *testing.T) {

	os.Setenv("AUCTION_INTERVAL", "10ms")

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Success", func(mt *mtest.T) {
		//Create auction
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		//Find auctions to complete
		first := mtest.CreateCursorResponse(1, "auctions.auctions", mtest.FirstBatch, bson.D{{"_id", "1"}})
		killCursors := mtest.CreateCursorResponse(0, "auctions.auctions", mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		//update auction
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

		ctx := context.Background()
		err := auctionRepository.CreateAuction(ctx, auction)
		assert.Nil(t, err)

		time.Sleep(50 * time.Millisecond)

		auctionRepository.CompleteAuction(ctx, 10*time.Millisecond)

		status, okauction := auctionStatusMap.GetAuctionStatus(auction.Id)
		assert.True(t, okauction)
		assert.Equal(t, auction_entity.Completed, status)
	})
}
