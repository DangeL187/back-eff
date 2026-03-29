package infra

import (
	"context"
	"errors"

	"github.com/DangeL187/erax"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"back/internal/features/subs/domain"
)

type CrudlRepo struct {
	db *gorm.DB
}

func NewCrudlRepo(db *gorm.DB) *CrudlRepo {
	return &CrudlRepo{db: db}
}

func (r *CrudlRepo) Create(ctx context.Context, s *domain.Subscription) error {
	err := r.db.WithContext(ctx).Create(s).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrSubscriptionAlreadyExists
		}

		err = erax.Wrap(err, "failed to create subscription")
		return erax.WithMeta(err, "layer", "DB")
	}

	return nil
}

func (r *CrudlRepo) GetByID(ctx context.Context, id int64) (*domain.Subscription, error) {
	var s domain.Subscription
	err := r.db.WithContext(ctx).First(&s, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrSubscriptionNotFound
	}
	if err != nil {
		err = erax.Wrap(err, "failed to get subscription by ID")
		return nil, erax.WithMeta(err, "layer", "DB")
	}

	return &s, nil
}

func (r *CrudlRepo) Update(ctx context.Context, s *domain.Subscription) error {
	err := r.db.WithContext(ctx).Save(s).Error
	if err != nil {
		err = erax.Wrap(err, "failed to update subscription")
		return erax.WithMeta(err, "layer", "DB")
	}

	return nil
}

func (r *CrudlRepo) DeleteByID(ctx context.Context, id int64) error {
	err := r.db.WithContext(ctx).Delete(&domain.Subscription{}, id).Error
	if err != nil {
		err = erax.Wrap(err, "failed to delete subscription by ID")
		return erax.WithMeta(err, "layer", "DB")
	}

	return nil
}

func (r *CrudlRepo) List(ctx context.Context, userID *uuid.UUID, service *string) ([]domain.Subscription, error) {
	var subs []domain.Subscription
	query := r.db.WithContext(ctx).Model(&domain.Subscription{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if service != nil {
		query = query.Where("service_name = ?", *service)
	}

	err := query.Find(&subs).Error
	if err != nil {
		err = erax.Wrap(err, "failed to find subscriptions")
		return nil, erax.WithMeta(err, "layer", "DB")
	}
	return subs, nil
}

func (r *CrudlRepo) Exists(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Subscription{}).
		Where("id = ?", id).Count(&count).Error
	if err != nil {
		err = erax.Wrap(err, "failed to count subscription")
		return false, erax.WithMeta(err, "layer", "DB")
	}
	return count > 0, nil
}
