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

type officerRepository struct {
	db *gorm.DB
}

func NewOfficerRepository(db *gorm.DB) contract.OfficerRepository {
	return &officerRepository{db: db}
}

func (r *officerRepository) FindIncomingOrders(ctx context.Context, pagination contract.PaginationFilter) ([]entity.LeasingContract, int64, error) {
	var contracts []entity.LeasingContract
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.LeasingContract{}).
		Where("status IN ?", []string{"draft", "pending", "approved"}).
		Preload("Customer").
		Preload("Motor")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pagination.Page - 1) * pagination.Limit
	err := query.Offset(offset).Limit(pagination.Limit).Order("created_at DESC").Find(&contracts).Error

	return contracts, total, err
}

func (r *officerRepository) FindTasksByRoleID(ctx context.Context, roleID int64, pagination contract.PaginationFilter) ([]entity.LeasingTask, int64, error) {
	var tasks []entity.LeasingTask
	var total int64

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

func (r *officerRepository) FindTaskByID(ctx context.Context, taskID int64) (*entity.LeasingTask, error) {
	var task entity.LeasingTask
	err := r.db.WithContext(ctx).First(&task, taskID).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *officerRepository) FindNextTask(ctx context.Context, contractID int64, currentSequence int16) (*entity.LeasingTask, error) {
	var task entity.LeasingTask
	err := r.db.WithContext(ctx).
		Where("contract_id = ? AND sequence_no > ?", contractID, currentSequence).
		Order("sequence_no ASC").
		First(&task).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil // No next task
	}
	return &task, nil
}

func (r *officerRepository) ProcessTaskAndUpdateNext(ctx context.Context, currentTask *entity.LeasingTask, nextTask *entity.LeasingTask, isFinal bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Update current task
		currentTask.Status = vo.StatusCompleted.String()
		currentTask.ActualEnddate = &now
		if err := tx.Save(currentTask).Error; err != nil {
			return fmt.Errorf("failed to update current task: %w", err)
		}

		// Activate next task
		if nextTask != nil {
			nextTask.Status = vo.StatusInProgress.String()
			nextTask.ActualStartdate = &now
			if err := tx.Save(nextTask).Error; err != nil {
				return fmt.Errorf("failed to update next task: %w", err)
			}
		}

		// Update contract status based on progress
		if isFinal {
			// All tasks done → contract is active
			if err := tx.Model(&entity.LeasingContract{}).
				Where("contract_id = ?", currentTask.ContractID).
				Update("status", vo.StatusActive.String()).Error; err != nil {
				return fmt.Errorf("failed to finalise contract status: %w", err)
			}
		} else {
			// Tasks still in progress → set contract to inprogress (if still draft)
			if err := tx.Model(&entity.LeasingContract{}).
				Where("contract_id = ? AND status = ?", currentTask.ContractID, vo.StatusDraft.String()).
				Update("status", vo.StatusApproved.String()).Error; err != nil {
				return fmt.Errorf("failed to update contract status: %w", err)
			}
		}

		return nil
	})
}

func (r *officerRepository) FindRoleIDByName(ctx context.Context, roleName string) (int64, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Where("role_name = ?", roleName).First(&role).Error
	if err != nil {
		return 0, err
	}
	return role.RoleID, nil
}
