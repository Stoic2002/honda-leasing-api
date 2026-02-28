package contract

import (
	"context"

	"honda-leasing-api/internal/domain/entity"
)

type FinanceRepository interface {
	FindPaymentSchedules(ctx context.Context, contractID int64) ([]entity.PaymentSchedule, error)
}
