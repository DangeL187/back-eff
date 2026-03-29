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

func TestDeleteEndpoint_Success(t *testing.T) {
	e, db := setupTestServer(t)

	sub := domain.Subscription{
		UserID:      uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
		ServiceName: "DeleteMe",
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		PricePerDay: 10,
		EndDate:     pgtype.Date{Valid: true, InfinityModifier: pgtype.Infinity},
	}
	assert.NoError(t, db.DB.Create(&sub).Error)

	req := httptest.NewRequest(http.MethodDelete, "/subscriptions/"+strconv.FormatInt(sub.ID, 10), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)

	var deleted domain.Subscription
	err := db.DB.First(&deleted, sub.ID).Error
	assert.Error(t, err)
}
