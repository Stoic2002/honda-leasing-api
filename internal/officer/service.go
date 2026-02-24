package officer

import (
	"context"
	"fmt"

	"honda-leasing-api/internal/domain"
	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/internal/domain/entity"
	"honda-leasing-api/internal/domain/vo"
	"honda-leasing-api/pkg/pagination"
)

// ProcessTaskInput is a clean service-level input struct (no HTTP tags).
type ProcessTaskInput struct {
	Notes      string
	Attributes map[string]string // Key: tasa_name, Value: string value or URL
}

type Service interface {
	GetIncomingOrders(ctx context.Context, pagination contract.PaginationFilter) ([]entity.LeasingContract, int64, error)
	GetMyTasks(ctx context.Context, userRoleName string, pagination contract.PaginationFilter) ([]entity.LeasingTask, int64, error)
	ProcessOrderTask(ctx context.Context, taskID int64, userRoleName string, req ProcessTaskInput) error

	RegisterCallFunction(name string, fn func(ctx context.Context, contractID int64) error)
}

type service struct {
	repo             contract.OfficerRepository
	functionRegistry map[string]func(ctx context.Context, contractID int64) error
}

func NewService(repo contract.OfficerRepository) Service {
	return &service{
		repo:             repo,
		functionRegistry: make(map[string]func(ctx context.Context, contractID int64) error),
	}
}

func (s *service) RegisterCallFunction(name string, fn func(ctx context.Context, contractID int64) error) {
	s.functionRegistry[name] = fn
}

func (s *service) GetIncomingOrders(ctx context.Context, pg contract.PaginationFilter) ([]entity.LeasingContract, int64, error) {
	pg.Page, pg.Limit = pagination.Normalize(pg.Page, pg.Limit)
	return s.repo.FindIncomingOrders(ctx, pg)
}

func (s *service) GetMyTasks(ctx context.Context, userRoleName string, pg contract.PaginationFilter) ([]entity.LeasingTask, int64, error) {
	pg.Page, pg.Limit = pagination.Normalize(pg.Page, pg.Limit)

	roleID, err := s.repo.FindRoleIDByName(ctx, userRoleName)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: invalid role %s", domain.ErrForbidden, userRoleName)
	}

	return s.repo.FindTasksByRoleID(ctx, roleID, pg)
}

func (s *service) ProcessOrderTask(ctx context.Context, taskID int64, userRoleName string, req ProcessTaskInput) error {
	// 1. Resolve the user's role_id from role name
	userRoleID, err := s.repo.FindRoleIDByName(ctx, userRoleName)
	if err != nil {
		return fmt.Errorf("%w: invalid role %s", domain.ErrForbidden, userRoleName)
	}

	// 2. Fetch current task
	currentTask, err := s.repo.FindTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("%w: task %d", domain.ErrNotFound, taskID)
	}

	// 3. Validate that the task belongs to the user's role
	if currentTask.RoleID != userRoleID {
		return fmt.Errorf("%w: this task is not assigned to your role", domain.ErrForbidden)
	}

	if currentTask.Status != vo.StatusInProgress.String() {
		return fmt.Errorf("%w: task %d is currently %s", domain.ErrInvalidInput, taskID, currentTask.Status)
	}

	// 4. See what the next task in the order's sequence is
	nextTask, err := s.repo.FindNextTask(ctx, currentTask.ContractID, currentTask.SequenceNo)
	if err != nil {
		return fmt.Errorf("%w: failed fetching next task", domain.ErrInternalServerError)
	}

	isFinal := nextTask == nil

	// 5. Process the transition transactionally
	err = s.repo.ProcessTaskAndUpdateNext(ctx, currentTask, nextTask, isFinal, req.Attributes)
	if err != nil {
		return fmt.Errorf("%w: failed processing logic", domain.ErrInternalServerError)
	}

	// 6. Execute Call Function dynamically if present
	if currentTask.CallFunction != nil && *currentTask.CallFunction != "" {
		fnName := *currentTask.CallFunction
		if fn, exists := s.functionRegistry[fnName]; exists {
			if err := fn(ctx, currentTask.ContractID); err != nil {
				return fmt.Errorf("failed executing call_function '%s': %w", fnName, err)
			}
		} else {
			// Fail softly or hard? We choose to fail hard if function is missing but declared
			return fmt.Errorf("fatal: call_function '%s' is not registered in service", fnName)
		}
	}

	return nil
}
