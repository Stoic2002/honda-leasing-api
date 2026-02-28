package postgres

import (
	"context"

	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/internal/domain/entity"

	"gorm.io/gorm"
)

type financeRepository struct {
	db *gorm.DB
}

func NewFinanceRepository(db *gorm.DB) contract.FinanceRepository {
	return &financeRepository{db: db}
}

func (r *financeRepository) FindPaymentSchedules(ctx context.Context, contractID int64) ([]entity.PaymentSchedule, error) {
	var schedules []entity.PaymentSchedule
	query := r.db.WithContext(ctx)
	if contractID > 0 {
		query = query.Where("contract_id = ?", contractID)
	}
	err := query.Order("angsuran_ke ASC").Find(&schedules).Error
	return schedules, err
}
