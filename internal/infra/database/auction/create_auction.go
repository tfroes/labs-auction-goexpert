package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection       *mongo.Collection
	AuctionStatusMap *AuctionStatusMap
}

func NewAuctionRepository(database *mongo.Database, auctionStatusMap *AuctionStatusMap) *AuctionRepository {
	return &AuctionRepository{
		Collection:       database.Collection("auctions"),
		AuctionStatusMap: auctionStatusMap,
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	go func() {
		select {
		case <-time.After(getAuctionInterval()):
			update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
			filter := bson.M{"_id": auctionEntityMongo.Id}

			_, err := ar.Collection.UpdateOne(ctx, filter, update)
			if err != nil {
				logger.Error("Error trying to update auction status to completed", err)
				return
			}

			ar.AuctionStatusMap.SetAuctionStatus(auctionEntityMongo.Id, auction_entity.Completed)
		}
	}()

	return nil
}
