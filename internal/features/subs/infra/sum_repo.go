package infra

import (
	"context"
	"time"

	"github.com/DangeL187/erax"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SumRepo struct {
	db *gorm.DB
}

func NewSumRepo(db *gorm.DB) *SumRepo {
	return &SumRepo{db: db}
}

// GetSubscriptionsCostForPeriod считает суммарную стоимость подписок за период [from, to]
// с опциональными фильтрами по userID и serviceName.
func (r *SumRepo) GetSubscriptionsCostForPeriod(ctx context.Context, from, to time.Time, userID *uuid.UUID, service *string) (int64, error) {
	var total int64

	// SQL-запрос:
	// 1) tsrange(start_date,end_date) && tsrange(from,to) -> пересечение подписки с периодом
	// 2) days_in_period = количество дней пересечения
	//      - если end_date = infinity -> берём конец периода to
	//      - иначе -> минимум(end_date, to)
	// 3) суммируем: price_per_day * days_in_period
	// 4) фильтры применяются только если переданы userID или service
	query := `
		SELECT COALESCE(
			SUM(
				price_per_day *
				(
					CASE
						WHEN end_date = 'infinity'::date THEN $2
						ELSE LEAST(end_date, $2)
					END
					- GREATEST(start_date, $1) + 1
				)
			), 0
		) AS total_cost
		FROM subscriptions
		WHERE period && tsrange($1::timestamp, $2::timestamp, '[]')
		  AND ($3::uuid IS NULL OR user_id = $3::uuid)
  		  AND ($4::text IS NULL OR service_name = $4::text);
	`

	err := r.db.WithContext(ctx).Raw(query, from, to, userID, service).Scan(&total).Error
	if err != nil {
		err = erax.Wrap(err, "failed to calculate subscription costs")
		zap.L().Error(erax.FormatToJSONString(err))
		return 0, err
	}

	return total, nil
}
