package test

import (
	"bytes"
	"encoding/json"
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

func TestUpdateEndpoint_Success(t *testing.T) {
	e, db := setupTestServer(t)

	sub := domain.Subscription{
		UserID:      uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
		ServiceName: "Old Service",
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		PricePerDay: 100,
		EndDate:     pgtype.Date{Valid: true, InfinityModifier: pgtype.Infinity},
	}
	assert.NoError(t, db.DB.Create(&sub).Error)

	reqBody := map[string]interface{}{
		"service_name": "Updated Service",
		"price":        200,
		"user_id":      sub.UserID.String(),
		"start_date":   "02-2025",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/subscriptions/"+strconv.FormatInt(sub.ID, 10), bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var updated domain.Subscription
	assert.NoError(t, db.DB.First(&updated, sub.ID).Error)
	assert.Equal(t, "Updated Service", updated.ServiceName)
	assert.Equal(t, int64(200), updated.PricePerDay)
}
