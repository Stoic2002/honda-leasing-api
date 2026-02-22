package postgres

import (
	"context"
	"fmt"

	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/internal/domain/entity"

	"gorm.io/gorm"
)

type leasingRepository struct {
	db *gorm.DB
}

func NewLeasingRepository(db *gorm.DB) contract.LeasingRepository {
	return &leasingRepository{db: db}
}

func (r *leasingRepository) CreateContract(ctx context.Context, contract *entity.LeasingContract, tasks []entity.LeasingTask) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Create Contract
		if err := tx.Create(contract).Error; err != nil {
			return fmt.Errorf("failed to create contract: %w", err)
		}

		// 2. Assign Contract ID to all tasks
		for i := range tasks {
			tasks[i].ContractID = contract.ContractID
		}

		// 3. Batch insert Tasks
		if len(tasks) > 0 {
			if err := tx.Create(&tasks).Error; err != nil {
				return fmt.Errorf("failed to create tasks: %w", err)
			}
		}

		return nil
	})
}

func (r *leasingRepository) GetContractProgress(ctx context.Context, contractID int64) ([]entity.LeasingTask, error) {
	var tasks []entity.LeasingTask
	err := r.db.WithContext(ctx).
		Where("contract_id = ?", contractID).
		Order("sequence_no ASC").
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *leasingRepository) GetTemplateTasks(ctx context.Context) ([]entity.TemplateTask, error) {
	var templates []entity.TemplateTask
	err := r.db.WithContext(ctx).
		Order("sequence_no ASC").
		Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}
func (r *leasingRepository) FindCustomerByUserID(ctx context.Context, userID int64) (*entity.Customer, error) {
	var customer entity.Customer
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&customer).Error
	if err != nil {
		return nil, fmt.Errorf("customer not found for user_id %d: %w", userID, err)
	}
	return &customer, nil
}

func (r *leasingRepository) FindContractsByUserID(ctx context.Context, userID int64, pagination contract.PaginationFilter) ([]entity.LeasingContract, int64, error) {
	// 1. Resolve customer_id from user_id
	customer, err := r.FindCustomerByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	// 2. Query contracts for this customer
	var contracts []entity.LeasingContract
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.LeasingContract{}).
		Where("customer_id = ?", customer.CustomerID).
		Preload("Motor")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pagination.Page - 1) * pagination.Limit
	err = query.Offset(offset).Limit(pagination.Limit).Order("created_at DESC").Find(&contracts).Error

	return contracts, total, err
}
