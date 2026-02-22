package postgres

import (
	"context"
	"fmt"
	"time"

	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/internal/domain/entity"
	"honda-leasing-api/internal/domain/vo"

	"gorm.io/gorm"
)

type deliveryRepository struct {
	db *gorm.DB
}

func NewDeliveryRepository(db *gorm.DB) contract.DeliveryRepository {
	return &deliveryRepository{db: db}
}

func (r *deliveryRepository) FindDeliveryOrders(ctx context.Context, roleID int64, pagination contract.PaginationFilter) ([]entity.LeasingContract, int64, error) {
	var contracts []entity.LeasingContract
	var total int64

	// Find contracts that have at least one inprogress task for this role
	subquery := r.db.Model(&entity.LeasingTask{}).
		Select("DISTINCT contract_id").
		Where("role_id = ? AND status = ?", roleID, "inprogress")

	query := r.db.WithContext(ctx).Model(&entity.LeasingContract{}).
		Where("contract_id IN (?)", subquery).
		Preload("Customer").
		Preload("Motor")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pagination.Page - 1) * pagination.Limit
	err := query.Offset(offset).Limit(pagination.Limit).Order("created_at DESC").Find(&contracts).Error

	return contracts, total, err
}

func (r *deliveryRepository) FindDeliveryTasks(ctx context.Context, roleID int64, pagination contract.PaginationFilter) ([]entity.LeasingTask, int64, error) {
	var tasks []entity.LeasingTask
	var total int64

	// Filter by role_id (SALES) and status inprogress
	query := r.db.WithContext(ctx).Model(&entity.LeasingTask{}).
		Where("role_id = ?", roleID).
		Where("status = ?", "inprogress")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pagination.Page - 1) * pagination.Limit
	err := query.Offset(offset).Limit(pagination.Limit).Order("created_at ASC").Find(&tasks).Error

	return tasks, total, err
}

func (r *deliveryRepository) FindTaskByID(ctx context.Context, taskID int64) (*entity.LeasingTask, error) {
	var task entity.LeasingTask
	err := r.db.WithContext(ctx).First(&task, taskID).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *deliveryRepository) GetContractByID(ctx context.Context, contractID int64) (*entity.LeasingContract, error) {
	var contract entity.LeasingContract
	err := r.db.WithContext(ctx).First(&contract, contractID).Error
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

func (r *deliveryRepository) CompleteDeliveryTask(ctx context.Context, currentTask *entity.LeasingTask, contract *entity.LeasingContract, asset *entity.MotorAsset) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// 1. Update current task to completed
		currentTask.Status = vo.StatusCompleted.String()
		currentTask.ActualEnddate = &now
		if err := tx.Save(currentTask).Error; err != nil {
			return fmt.Errorf("failed completing task: %w", err)
		}

		// 2. Insert Asset Proof (if provided)
		if asset != nil {
			if err := tx.Create(asset).Error; err != nil {
				return fmt.Errorf("failed saving delivery proof asset: %w", err)
			}
		}

		// 3. Find and activate the next task
		var nextTask entity.LeasingTask
		err := tx.Where("contract_id = ? AND sequence_no > ?", currentTask.ContractID, currentTask.SequenceNo).
			Order("sequence_no ASC").
			First(&nextTask).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed finding next task: %w", err)
		}

		if err == nil {
			// There IS a next task → activate it
			nextTask.Status = vo.StatusInProgress.String()
			nextTask.ActualStartdate = &now
			if err := tx.Save(&nextTask).Error; err != nil {
				return fmt.Errorf("failed activating next task: %w", err)
			}

			// Update contract to inprogress (if still draft)
			if contract.Status == vo.StatusDraft.String() {
				contract.Status = vo.StatusApproved.String()
				if err := tx.Save(contract).Error; err != nil {
					return fmt.Errorf("failed updating contract status: %w", err)
				}
			}
		} else {
			// No next task → this is the final step, set contract to active
			contract.Status = vo.StatusActive.String()
			contract.TanggalMulaiCicil = &now
			if err := tx.Save(contract).Error; err != nil {
				return fmt.Errorf("failed finalising contract status: %w", err)
			}
		}

		return nil
	})
}

func (r *deliveryRepository) FindRoleIDByName(ctx context.Context, roleName string) (int64, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Where("role_name = ?", roleName).First(&role).Error
	if err != nil {
		return 0, err
	}
	return role.RoleID, nil
}
