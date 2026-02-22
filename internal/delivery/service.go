package delivery

import (
	"context"
	"fmt"

	"honda-leasing-api/internal/domain"
	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/internal/domain/entity"
	"honda-leasing-api/internal/domain/vo"
	"honda-leasing-api/pkg/pagination"
)

type CompleteDeliveryRequest struct {
	FileName string
	FileSize float64
	FileType string
	FileURL  string
}

type Service interface {
	GetDeliveryOrders(ctx context.Context, userRoleName string, pagination contract.PaginationFilter) ([]entity.LeasingContract, int64, error)
	GetDeliveryTasks(ctx context.Context, userRoleName string, pagination contract.PaginationFilter) ([]entity.LeasingTask, int64, error)
	CompleteDelivery(ctx context.Context, taskID int64, userRoleName string, req CompleteDeliveryRequest) error
}

type service struct {
	repo contract.DeliveryRepository
}

func NewService(repo contract.DeliveryRepository) Service {
	return &service{repo: repo}
}

func (s *service) GetDeliveryOrders(ctx context.Context, userRoleName string, pg contract.PaginationFilter) ([]entity.LeasingContract, int64, error) {
	pg.Page, pg.Limit = pagination.Normalize(pg.Page, pg.Limit)

	roleID, err := s.repo.FindRoleIDByName(ctx, userRoleName)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: invalid role %s", domain.ErrForbidden, userRoleName)
	}

	return s.repo.FindDeliveryOrders(ctx, roleID, pg)
}

func (s *service) GetDeliveryTasks(ctx context.Context, userRoleName string, pg contract.PaginationFilter) ([]entity.LeasingTask, int64, error) {
	pg.Page, pg.Limit = pagination.Normalize(pg.Page, pg.Limit)

	// Resolve role_id from role name
	roleID, err := s.repo.FindRoleIDByName(ctx, userRoleName)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: invalid role %s", domain.ErrForbidden, userRoleName)
	}

	return s.repo.FindDeliveryTasks(ctx, roleID, pg)
}

func (s *service) CompleteDelivery(ctx context.Context, taskID int64, userRoleName string, req CompleteDeliveryRequest) error {
	// 1. Resolve the user's role_id
	userRoleID, err := s.repo.FindRoleIDByName(ctx, userRoleName)
	if err != nil {
		return fmt.Errorf("%w: invalid role %s", domain.ErrForbidden, userRoleName)
	}

	// 2. Validate task
	currentTask, err := s.repo.FindTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("%w: task %d", domain.ErrNotFound, taskID)
	}

	// 3. Validate that the task belongs to the user's role
	if currentTask.RoleID != userRoleID {
		return fmt.Errorf("%w: this task is not assigned to your role", domain.ErrForbidden)
	}

	if currentTask.Status != vo.StatusInProgress.String() {
		return fmt.Errorf("%w: delivery task %d is currently %s", domain.ErrInvalidInput, taskID, currentTask.Status)
	}

	// 4. Fetch associated contract to get motor ID
	c, err := s.repo.GetContractByID(ctx, currentTask.ContractID)
	if err != nil {
		return fmt.Errorf("%w: contract for task not found", domain.ErrNotFound)
	}

	// 5. Prepare asset entity for proof
	var asset *entity.MotorAsset
	if req.FileURL != "" {
		asset = &entity.MotorAsset{
			FileName:    req.FileName,
			FileSize:    req.FileSize,
			FileType:    req.FileType,
			FileURL:     req.FileURL,
			MoasMotorID: c.MotorID,
		}
	}

	// 6. Fire the complete transaction
	return s.repo.CompleteDeliveryTask(ctx, currentTask, c, asset)
}
