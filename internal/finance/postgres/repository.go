package postgres

import (
	"context"
	"errors"

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

func (r *financeRepository) FindScheduleByID(ctx context.Context, scheduleID int64) (*entity.PaymentSchedule, error) {
	var schedule entity.PaymentSchedule
	err := r.db.WithContext(ctx).First(&schedule, scheduleID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &schedule, nil
}

func (r *financeRepository) CreatePaymentAndUpdateSchedule(ctx context.Context, payment *entity.Payment, schedule *entity.PaymentSchedule) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(payment).Error; err != nil {
			return err
		}
		if err := tx.Save(schedule).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *financeRepository) CreatePaymentSchedules(ctx context.Context, schedules []entity.PaymentSchedule) error {
	return r.db.WithContext(ctx).Create(&schedules).Error
}

func (r *financeRepository) FindContractByID(ctx context.Context, contractID int64) (*entity.LeasingContract, error) {
	var contract entity.LeasingContract
	err := r.db.WithContext(ctx).First(&contract, contractID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &contract, nil
}
