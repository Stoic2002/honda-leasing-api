package catalog

import (
	"context"

	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/internal/domain/entity"
	"honda-leasing-api/pkg/pagination"
)

type Service interface {
	GetMotors(ctx context.Context, filter contract.CatalogFilter, pagination contract.PaginationFilter) ([]entity.Motor, int64, error)
	GetMotorByID(ctx context.Context, id int64) (*entity.Motor, error)
	GetLeasingProducts(ctx context.Context) ([]entity.LeasingProduct, error)
}

type service struct {
	repo contract.CatalogRepository
}

func NewService(repo contract.CatalogRepository) Service {
	return &service{repo: repo}
}

func (s *service) GetMotors(ctx context.Context, filter contract.CatalogFilter, pg contract.PaginationFilter) ([]entity.Motor, int64, error) {
	pg.Page, pg.Limit = pagination.Normalize(pg.Page, pg.Limit)
	return s.repo.FindAllMotors(ctx, filter, pg)
}

func (s *service) GetMotorByID(ctx context.Context, id int64) (*entity.Motor, error) {
	return s.repo.FindMotorByID(ctx, id)
}

func (s *service) GetLeasingProducts(ctx context.Context) ([]entity.LeasingProduct, error) {
	return s.repo.GetLeasingProducts(ctx)
}
