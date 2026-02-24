package contract

import (
	"context"

	"honda-leasing-api/internal/domain/entity"
)

type FinanceRepository interface {
	FindPaymentSchedules(ctx context.Context, contractID int64) ([]entity.PaymentSchedule, error)
	FindScheduleByID(ctx context.Context, scheduleID int64) (*entity.PaymentSchedule, error)
	CreatePaymentAndUpdateSchedule(ctx context.Context, payment *entity.Payment, schedule *entity.PaymentSchedule) error
}
