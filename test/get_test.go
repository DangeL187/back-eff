package test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	"back/internal/features/subs/domain"
)

func TestGetEndpoint_Success(t *testing.T) {
	e, db := setupTestServer(t)

	// Создаем запись напрямую
	sub := domain.Subscription{
		UserID:      uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
		ServiceName: "YouTube Premium",
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		PricePerDay: 50,
		EndDate:     pgtype.Date{Valid: true, InfinityModifier: pgtype.Infinity},
	}
	assert.NoError(t, db.DB.Create(&sub).Error)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/"+strconv.FormatInt(sub.ID, 10), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetEndpoint_NotFound(t *testing.T) {
	e, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/999999", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetEndpoint_InvalidID(t *testing.T) {
	e, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/abc", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
