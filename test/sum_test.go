package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	"back/internal/features/subs/domain"
)

func TestGetSubscriptionCostsForPeriod_Success(t *testing.T) {
	e, db := setupTestServer(t)

	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	sub := domain.Subscription{
		UserID:      userID,
		ServiceName: "Netflix",
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		PricePerDay: 100,
		EndDate:     pgtype.Date{Valid: true, InfinityModifier: pgtype.Infinity},
	}
	assert.NoError(t, db.DB.Create(&sub).Error)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/costs?user_id="+userID.String()+"&from=01-01-2025&to=31-12-2025", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result int64
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)

	// 365 дней * 100 = 36500
	assert.Equal(t, int64(36500), result)
}

func TestGetSubscriptionCostsForPeriod_WithEndDate(t *testing.T) {
	e, db := setupTestServer(t)

	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cbb")
	sub := domain.Subscription{
		UserID:      userID,
		ServiceName: "HBO Max",
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     pgtype.Date{Time: time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC), Valid: true, InfinityModifier: pgtype.Finite},
		PricePerDay: 50,
	}
	assert.NoError(t, db.DB.Create(&sub).Error)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/costs?user_id="+userID.String()+"&from=01-01-2025&to=31-12-2025", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result int64
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)

	// 181 день * 50 = 9050
	assert.Equal(t, int64(9050), result)
}

func TestGetSubscriptionCostsForPeriod_NoFilters(t *testing.T) {
	e, db := setupTestServer(t)

	// Две подписки
	sub1 := domain.Subscription{
		UserID:      uuid.New(),
		ServiceName: "Netflix",
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     pgtype.Date{Valid: true, InfinityModifier: pgtype.Infinity},
		PricePerDay: 100,
	}
	sub2 := domain.Subscription{
		UserID:      uuid.New(),
		ServiceName: "HBO Max",
		StartDate:   time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     pgtype.Date{Time: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC), Valid: true, InfinityModifier: pgtype.Finite},
		PricePerDay: 50,
	}
	assert.NoError(t, db.DB.Create(&sub1).Error)
	assert.NoError(t, db.DB.Create(&sub2).Error)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/costs?from=01-01-2025&to=31-12-2025", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result int64
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)

	// sub1: 365 дней * 100 = 36500
	// sub2: февраль 2025 = 28 дней * 50 = 1400
	assert.Equal(t, int64(36500+1400), result)
}

func TestGetSubscriptionCostsForPeriod_FilterUserAndService(t *testing.T) {
	e, db := setupTestServer(t)

	userID := uuid.New()
	// подписка, подходящая под фильтр
	sub := domain.Subscription{
		UserID:      userID,
		ServiceName: "Disney+",
		StartDate:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     pgtype.Date{Valid: true, InfinityModifier: pgtype.Infinity},
		PricePerDay: 30,
	}
	// подписка другого пользователя
	subOther := domain.Subscription{
		UserID:      uuid.New(),
		ServiceName: "Disney+",
		StartDate:   time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     pgtype.Date{Valid: true, InfinityModifier: pgtype.Infinity},
		PricePerDay: 30,
	}

	assert.NoError(t, db.DB.Create(&sub).Error)
	assert.NoError(t, db.DB.Create(&subOther).Error)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/costs?user_id="+userID.String()+"&service=Disney%2B&from=01-03-2025&to=31-03-2025", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var result int64
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)

	// 31 день * 30 = 930
	assert.Equal(t, int64(930), result)
}

func TestGetSubscriptionCostsForPeriod_InvalidDate(t *testing.T) {
	e, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/costs?from=01-01-2025&to=notadate", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
