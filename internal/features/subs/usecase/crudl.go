package usecase

import (
	"back/internal/features/subs/domain"
	"back/internal/features/subs/dto"
	"context"

	"github.com/DangeL187/erax"
	"github.com/google/uuid"
)

type CrudlUseCase struct {
	crudlRepo domain.CrudlRepo
}

func NewCrudlUseCase(crudlRepo domain.CrudlRepo) *CrudlUseCase {
	return &CrudlUseCase{crudlRepo: crudlRepo}
}

func (u *CrudlUseCase) Create(ctx context.Context, s dto.SubscriptionDTO) error {
	sub := &domain.Subscription{
		ID:          s.ID,
		UserID:      s.UserID,
		ServiceName: s.ServiceName,
		StartDate:   s.StartDate,
		EndDate:     s.EndDate,
		PricePerDay: s.PricePerDay,
	}

	err := u.crudlRepo.Create(ctx, sub)
	if err != nil {
		return erax.Wrap(err, "failed to create subscription")
	}

	return nil
}

func (u *CrudlUseCase) GetByID(ctx context.Context, id int64) (dto.SubscriptionDTO, error) {
	subFromRepo, err := u.crudlRepo.GetByID(ctx, id)
	if err != nil {
		return dto.SubscriptionDTO{}, erax.Wrap(err, "failed to get subscription by ID")
	}

	sub := dto.SubscriptionDTO{
		ID:          subFromRepo.ID,
		UserID:      subFromRepo.UserID,
		ServiceName: subFromRepo.ServiceName,
		StartDate:   subFromRepo.StartDate,
		EndDate:     subFromRepo.EndDate,
		PricePerDay: subFromRepo.PricePerDay,
	}

	return sub, nil
}

func (u *CrudlUseCase) Update(ctx context.Context, s dto.SubscriptionDTO) error {
	sub := &domain.Subscription{
		ID:          s.ID,
		UserID:      s.UserID,
		ServiceName: s.ServiceName,
		StartDate:   s.StartDate,
		EndDate:     s.EndDate,
		PricePerDay: s.PricePerDay,
	}

	exists, err := u.crudlRepo.Exists(ctx, sub.ID)
	if err != nil {
		return erax.Wrap(err, "failed to check subscription existence")
	}
	if !exists {
		return domain.ErrSubscriptionNotFound
	}

	err = u.crudlRepo.Update(ctx, sub)
	if err != nil {
		return erax.Wrap(err, "failed to update subscription")
	}

	return nil
}

func (u *CrudlUseCase) DeleteByID(ctx context.Context, id int64) error {
	err := u.crudlRepo.DeleteByID(ctx, id)
	if err != nil {
		return erax.Wrap(err, "failed to delete subscription by ID")
	}

	return nil
}

func (u *CrudlUseCase) List(ctx context.Context, userID *uuid.UUID, service *string) ([]dto.SubscriptionDTO, error) {
	subsFromRepo, err := u.crudlRepo.List(ctx, userID, service)
	if err != nil {
		return nil, erax.Wrap(err, "failed to list subscriptions")
	}

	var subs []dto.SubscriptionDTO
	for _, sub := range subsFromRepo {
		subs = append(subs, dto.SubscriptionDTO{
			ID:          sub.ID,
			UserID:      sub.UserID,
			ServiceName: sub.ServiceName,
			StartDate:   sub.StartDate,
			EndDate:     sub.EndDate,
			PricePerDay: sub.PricePerDay,
		})
	}

	return subs, nil
}
