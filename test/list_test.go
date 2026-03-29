package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	"back/internal/features/subs/domain"
)

func TestListEndpoint_FilterByUserID(t *testing.T) {
	e, db := setupTestServer(t)

	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	subs := []domain.Subscription{
		{UserID: userID, ServiceName: "A", StartDate: time.Now(), PricePerDay: 100, EndDate: pgtype.Date{Valid: true, InfinityModifier: pgtype.Infinity}},
		{UserID: userID, ServiceName: "B", StartDate: time.Now(), PricePerDay: 200, EndDate: pgtype.Date{Valid: true, InfinityModifier: pgtype.Infinity}},
	}
	for _, s := range subs {
		assert.NoError(t, db.DB.Create(&s).Error)
	}

	req := httptest.NewRequest(http.MethodGet, "/subscriptions?user_id="+userID.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
