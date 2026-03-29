package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	"back/internal/features/subs/domain"
)

func TestCreateEndpoint_WithoutEndDate(t *testing.T) {
	e, db := setupTestServer(t)

	reqBody := map[string]interface{}{
		"service_name": "Yandex Plus",
		"price":        400,
		"user_id":      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		"start_date":   "07-2025",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var sub domain.Subscription
	err := db.DB.
		Where("service_name = ? AND user_id = ?", "Yandex Plus", "60601fee-2bf1-4721-ae6f-7636e79a0cba").
		First(&sub).Error
	assert.NoError(t, err, "subscription should exist in DB")
	assert.Equal(t, int64(400), sub.PricePerDay)
	assert.Equal(t, "Yandex Plus", sub.ServiceName)
	assert.Equal(t, "60601fee-2bf1-4721-ae6f-7636e79a0cba", sub.UserID.String())
	assert.Equal(t, time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC), sub.StartDate)
	assert.Equal(t, pgtype.Infinity, sub.EndDate.InfinityModifier)
}

func TestCreateEndpoint_WithEndDate(t *testing.T) {
	e, db := setupTestServer(t)

	reqBody := map[string]interface{}{
		"service_name": "Netflix",
		"price":        300,
		"user_id":      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		"start_date":   "07-2025",
		"end_date":     "12-2025",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var sub domain.Subscription
	err := db.DB.Where("service_name = ? AND user_id = ?", "Netflix", "60601fee-2bf1-4721-ae6f-7636e79a0cba").
		First(&sub).Error
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC), sub.EndDate.Time)
	assert.Equal(t, pgtype.Finite, sub.EndDate.InfinityModifier)
}

func TestCreateEndpoint_InvalidBody(t *testing.T) {
	e, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewReader([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateEndpoint_InvalidUserID(t *testing.T) {
	e, _ := setupTestServer(t)

	reqBody := map[string]interface{}{
		"service_name": "Spotify",
		"price":        100,
		"user_id":      "not-a-uuid",
		"start_date":   "07-2025",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
