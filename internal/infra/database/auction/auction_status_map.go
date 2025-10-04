package auction

import (
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"sync"
	"time"
)

type AuctionStatusMap struct {
	auctionInterval       time.Duration
	auctionStatusMap      map[string]auction_entity.AuctionStatus
	auctionEndTimeMap     map[string]time.Time
	auctionStatusMapMutex *sync.Mutex
	auctionEndTimeMutex   *sync.Mutex
}

func NewAuctionStatusMap() *AuctionStatusMap {
	return &AuctionStatusMap{
		auctionInterval:       getAuctionInterval(),
		auctionStatusMap:      make(map[string]auction_entity.AuctionStatus),
		auctionEndTimeMap:     make(map[string]time.Time),
		auctionStatusMapMutex: &sync.Mutex{},
		auctionEndTimeMutex:   &sync.Mutex{},
	}
}

func (asm *AuctionStatusMap) GetAuctionStatus(auctionId string) (auction_entity.AuctionStatus, bool) {
	asm.auctionStatusMapMutex.Lock()
	defer asm.auctionStatusMapMutex.Unlock()

	auctionStatus, ok := asm.auctionStatusMap[auctionId]
	return auctionStatus, ok
}

func (asm *AuctionStatusMap) SetAuctionStatus(auctionId string, auctionStatus auction_entity.AuctionStatus) {
	asm.auctionStatusMapMutex.Lock()
	defer asm.auctionStatusMapMutex.Unlock()

	asm.auctionStatusMap[auctionId] = auctionStatus
}

func (asm *AuctionStatusMap) GetAuctionEndTime(auctionId string) (time.Time, bool) {
	asm.auctionEndTimeMutex.Lock()
	defer asm.auctionEndTimeMutex.Unlock()

	auctionEndTime, ok := asm.auctionEndTimeMap[auctionId]
	return auctionEndTime, ok
}

func (asm *AuctionStatusMap) SetAuctionEndTime(auctionId string, auctionCreatedTime time.Time) {
	asm.auctionEndTimeMutex.Lock()
	defer asm.auctionEndTimeMutex.Unlock()

	asm.auctionEndTimeMap[auctionId] = auctionCreatedTime.Add(asm.auctionInterval)
}

func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 5
	}

	return duration
}
