package usecase

import (
	"context"
	"time"

	"github.com/DangeL187/erax"
	"github.com/google/uuid"

	"back/internal/features/subs/domain"
)

type SumUseCase struct {
	sumRepo domain.SumRepo
}

func NewSumUseCase(sumRepo domain.SumRepo) *SumUseCase {
	return &SumUseCase{sumRepo: sumRepo}
}

func (u *SumUseCase) GetSubscriptionsCostForPeriod(ctx context.Context, from, to time.Time, userID *uuid.UUID, service *string) (int64, error) {
	total, err := u.sumRepo.GetSubscriptionsCostForPeriod(ctx, from, to, userID, service)
	if err != nil {
		return 0, erax.Wrap(err, "failed to get subscriptions cost for period")
	}

	return total, nil
}
