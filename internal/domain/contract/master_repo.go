package contract

import (
	"context"

	"honda-leasing-api/internal/domain/entity"
)

type MasterRepository interface {
	FindProvinces(ctx context.Context) ([]entity.Province, error)
	FindKabupatensByProvID(ctx context.Context, provID int64) ([]entity.Kabupaten, error)
	FindKecamatansByKabID(ctx context.Context, kabID int64) ([]entity.Kecamatan, error)
	FindKelurahansByKecID(ctx context.Context, kecID int64) ([]entity.Kelurahan, error)
}
