package service

import (
	"fmt"
	"time"

	"github.com/DevPooks/payflow/backend/internal/domain"
	"github.com/google/uuid"
)

func BuildInstallments(planID string, total int64, count int, firstDue time.Time) ([]domain.Installment, error) {
	if total <= 0 {
		return nil, fmt.Errorf("%w: order amount must be positive", domain.ErrInvalidInput)
	}
	if count < 2 || count > 4 {
		return nil, domain.ErrInvalidInstallmentCount
	}
	if total < int64(count) {
		return nil, fmt.Errorf("%w: total must provide at least one minor unit per installment", domain.ErrInvalidInput)
	}
	base, remainder := total/int64(count), total%int64(count)
	items := make([]domain.Installment, count)
	for index := 0; index < count; index++ {
		amount := base
		if int64(index) < remainder {
			amount++
		}
		items[index] = domain.Installment{
			ID: uuid.NewString(), PlanID: planID, Sequence: index + 1,
			Amount: amount, Status: "pending", DueAt: addMonthsClamped(firstDue, index),
		}
	}
	return items, nil
}

func addMonthsClamped(value time.Time, months int) time.Time {
	firstOfMonth := time.Date(value.Year(), value.Month()+time.Month(months), 1, value.Hour(), value.Minute(), value.Second(), value.Nanosecond(), value.Location())
	lastDay := time.Date(firstOfMonth.Year(), firstOfMonth.Month()+1, 0, value.Hour(), value.Minute(), value.Second(), value.Nanosecond(), value.Location()).Day()
	day := value.Day()
	if day > lastDay {
		day = lastDay
	}
	return time.Date(firstOfMonth.Year(), firstOfMonth.Month(), day, value.Hour(), value.Minute(), value.Second(), value.Nanosecond(), value.Location())
}
