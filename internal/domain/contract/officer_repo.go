package contract

import (
	"context"

	"honda-leasing-api/internal/domain/entity"
)

type OfficerRepository interface {
	FindIncomingOrders(ctx context.Context, pagination PaginationFilter) ([]entity.LeasingContract, int64, error)
	FindTasksByRoleID(ctx context.Context, roleID int64, pagination PaginationFilter) ([]entity.LeasingTask, int64, error)
	FindTaskByID(ctx context.Context, taskID int64) (*entity.LeasingTask, error)
	FindNextTask(ctx context.Context, contractID int64, currentSequence int16) (*entity.LeasingTask, error)
	ProcessTaskAndUpdateNext(ctx context.Context, currentTask *entity.LeasingTask, nextTask *entity.LeasingTask, isFinal bool, attributes map[string]string) error
	FindRoleIDByName(ctx context.Context, roleName string) (int64, error)
}
