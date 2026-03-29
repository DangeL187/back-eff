package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CrudlRepo interface {
	Create(ctx context.Context, s *Subscription) error
	GetByID(ctx context.Context, id int64) (*Subscription, error)
	Update(ctx context.Context, s *Subscription) error
	DeleteByID(ctx context.Context, id int64) error
	List(ctx context.Context, userID *uuid.UUID, service *string) ([]Subscription, error)
	Exists(ctx context.Context, id int64) (bool, error)
}

type SumRepo interface {
	GetSubscriptionsCostForPeriod(ctx context.Context, from, to time.Time, userID *uuid.UUID, service *string) (int64, error)
}
