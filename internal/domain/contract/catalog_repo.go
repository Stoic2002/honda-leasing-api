package contract

import (
	"context"

	"honda-leasing-api/internal/domain/entity"
)

type CatalogFilter struct {
	Search    string
	MotorType string
}

type PaginationFilter struct {
	Page  int
	Limit int
}

type CatalogRepository interface {
	FindAllMotors(ctx context.Context, filter CatalogFilter, pagination PaginationFilter) ([]entity.Motor, int64, error)
	FindMotorByID(ctx context.Context, id int64) (*entity.Motor, error)
	GetLeasingProducts(ctx context.Context) ([]entity.LeasingProduct, error)
}
