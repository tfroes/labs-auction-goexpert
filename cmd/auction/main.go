package main

import (
	"context"
	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/infra/api/web/controller/auction_controller"
	"fullcycle-auction_go/internal/infra/api/web/controller/bid_controller"
	"fullcycle-auction_go/internal/infra/api/web/controller/user_controller"
	"fullcycle-auction_go/internal/infra/database/auction"
	"fullcycle-auction_go/internal/infra/database/bid"
	"fullcycle-auction_go/internal/infra/database/user"
	"fullcycle-auction_go/internal/usecase/auction_usecase"
	"fullcycle-auction_go/internal/usecase/bid_usecase"
	"fullcycle-auction_go/internal/usecase/user_usecase"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	userController, bidController, auctionsController, auctionUseCase := initDependencies(databaseConnection)
	completeauctionDuration := getCompleteAuctionInterval()

	go func() {
		for {
			auctionUseCase.CompleteAuction(ctx, getAuctionInterval())
			time.Sleep(completeauctionDuration)
		}
	}()

	router := gin.Default()
	//Auction
	router.GET("/auction", auctionsController.FindAuctions)
	router.GET("/auction/:auctionId", auctionsController.FindAuctionById)
	router.POST("/auction", auctionsController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionsController.FindWinningBidByAuctionId)
	//BID
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidByAuctionId)
	//User
	router.GET("/user/:userId", userController.FindUserById)
	router.POST("/user", userController.CreateUser)

	router.Run(":8080")
}

func initDependencies(database *mongo.Database) (
	*user_controller.UserController,
	*bid_controller.BidController,
	*auction_controller.AuctionController,
	auction_usecase.AuctionUseCaseInterface) {

	auctionStatusMap := auction.NewAuctionStatusMap()

	auctionRepository := auction.NewAuctionRepository(database, auctionStatusMap)
	bidRepository := bid.NewBidRepository(database, auctionRepository, auctionStatusMap)
	userRepository := user.NewUserRepository(database)

	auctionUseCase := auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository)

	userController := user_controller.NewUserController(
		user_usecase.NewUserUseCase(userRepository))
	auctionController := auction_controller.NewAuctionController(auctionUseCase)
	bidController := bid_controller.NewBidController(bid_usecase.NewBidUseCase(bidRepository))

	return userController, bidController, auctionController, auctionUseCase
}

func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 5
	}

	return duration
}

func getCompleteAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("COMPLETE_AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 5
	}

	return duration
}
