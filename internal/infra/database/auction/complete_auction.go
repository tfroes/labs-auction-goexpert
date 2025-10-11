package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (ar *AuctionRepository) CompleteAuction(
	ctx context.Context,
	durationCompleted time.Duration) *internal_error.InternalError {

	timeThreshould := time.Now().Add(-durationCompleted).Unix()

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"status", 0}},
				bson.D{{"timestamp", bson.D{{"$lt", timeThreshould}}}},
			},
		},
	}

	opts := options.Find().SetProjection(bson.M{"_id": 1})

	cursor, err := ar.Collection.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Error finding auctions", err)
		return internal_error.NewInternalServerError("Error finding auctions")
	}
	defer cursor.Close(ctx)

	var auctionsMongo []AuctionEntityMongo
	if err := cursor.All(ctx, &auctionsMongo); err != nil {
		logger.Error("Error decoding auctions", err)
		return internal_error.NewInternalServerError("Error decoding auctions")
	}

	for _, a := range auctionsMongo {
		update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
		filter := bson.M{"_id": a.Id}

		_, err := ar.Collection.UpdateOne(ctx, filter, update)
		if err != nil {
			logger.Error("Error trying to update auction status to completed", err)
			return internal_error.NewInternalServerError("Error update auctions status to completed")
		}

		ar.AuctionStatusMap.SetAuctionStatus(a.Id, auction_entity.Completed)
	}

	return nil
}
