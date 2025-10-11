package auction_usecase

import (
	"context"
	"fullcycle-auction_go/internal/internal_error"
	"time"
)

func (au *AuctionUseCase) CompleteAuction(
	ctx context.Context,
	durationCompleted time.Duration) *internal_error.InternalError {

	if err := au.auctionRepositoryInterface.CompleteAuction(
		ctx, durationCompleted); err != nil {
		return err
	}

	return nil
}
