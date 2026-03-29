package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SubscriptionDTO struct {
	ID          int64
	UserID      uuid.UUID
	ServiceName string
	StartDate   time.Time
	EndDate     pgtype.Date
	PricePerDay int64
}

type SubscriptionDTOSwagger struct {
	UserID      uuid.UUID `json:"user_id"`
	ServiceName string    `json:"service_name"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	PricePerDay float64   `json:"price_per_day"`
}

type SubscriptionJSON struct {
	UserID      string `json:"user_id"`
	ServiceName string `json:"service_name"`
	StartDate   string `json:"start_date"`         // MM-YYYY
	EndDate     string `json:"end_date,omitempty"` // MM-YYYY
	PricePerDay int64  `json:"price"`
}
