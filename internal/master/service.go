package master

import (
	"context"

	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/internal/domain/entity"
)

type Service interface {
	GetProvinces(ctx context.Context) ([]entity.Province, error)
	GetKabupatens(ctx context.Context, provID int64) ([]entity.Kabupaten, error)
	GetKecamatans(ctx context.Context, kabID int64) ([]entity.Kecamatan, error)
	GetKelurahans(ctx context.Context, kecID int64) ([]entity.Kelurahan, error)
}

type service struct {
	repo contract.MasterRepository
}

func NewService(repo contract.MasterRepository) Service {
	return &service{repo: repo}
}

func (s *service) GetProvinces(ctx context.Context) ([]entity.Province, error) {
	return s.repo.FindProvinces(ctx)
}

func (s *service) GetKabupatens(ctx context.Context, provID int64) ([]entity.Kabupaten, error) {
	return s.repo.FindKabupatensByProvID(ctx, provID)
}

func (s *service) GetKecamatans(ctx context.Context, kabID int64) ([]entity.Kecamatan, error) {
	return s.repo.FindKecamatansByKabID(ctx, kabID)
}

func (s *service) GetKelurahans(ctx context.Context, kecID int64) ([]entity.Kelurahan, error) {
	return s.repo.FindKelurahansByKecID(ctx, kecID)
}
