package contract

import (
	"context"

	"honda-leasing-api/internal/domain/entity"
)

type DeliveryRepository interface {
	FindDeliveryOrders(ctx context.Context, roleID int64, pagination PaginationFilter) ([]entity.LeasingContract, int64, error)
	FindDeliveryTasks(ctx context.Context, roleID int64, pagination PaginationFilter) ([]entity.LeasingTask, int64, error)
	FindTaskByID(ctx context.Context, taskID int64) (*entity.LeasingTask, error)
	GetContractByID(ctx context.Context, contractID int64) (*entity.LeasingContract, error)
	CompleteDeliveryTask(ctx context.Context, currentTask *entity.LeasingTask, contract *entity.LeasingContract, asset *entity.MotorAsset) error
	FindRoleIDByName(ctx context.Context, roleName string) (int64, error)
}
