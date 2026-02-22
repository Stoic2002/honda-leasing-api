package contract

import (
	"context"

	"honda-leasing-api/internal/domain/entity"
)

type LeasingRepository interface {
	CreateContract(ctx context.Context, contract *entity.LeasingContract, tasks []entity.LeasingTask) error
	GetContractProgress(ctx context.Context, contractID int64) ([]entity.LeasingTask, error)
	GetTemplateTasks(ctx context.Context) ([]entity.TemplateTask, error)
	FindCustomerByUserID(ctx context.Context, userID int64) (*entity.Customer, error)
	FindContractsByUserID(ctx context.Context, userID int64, pagination PaginationFilter) ([]entity.LeasingContract, int64, error)
}
