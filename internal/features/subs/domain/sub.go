package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Subscription struct {
	ID          int64       `gorm:"primaryKey;autoIncrement"`
	UserID      uuid.UUID   `gorm:"type:uuid;index:idx_subscriptions_user_service,priority:1;not null"`
	ServiceName string      `gorm:"index:idx_subscriptions_user_service,priority:2;not null"`
	StartDate   time.Time   `gorm:"not null"`
	EndDate     pgtype.Date `gorm:"not null;default:'infinity'::date"`
	PricePerDay int64       `gorm:"not null"`
	Period      string      `gorm:"->;type:tsrange;generatedAlways"`
}
