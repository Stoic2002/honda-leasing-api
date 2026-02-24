package postgres

import (
	"context"

	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/internal/domain/entity"

	"gorm.io/gorm"
)

type masterRepository struct {
	db *gorm.DB
}

func NewMasterRepository(db *gorm.DB) contract.MasterRepository {
	return &masterRepository{db: db}
}

func (r *masterRepository) FindProvinces(ctx context.Context) ([]entity.Province, error) {
	var provinces []entity.Province
	err := r.db.WithContext(ctx).Find(&provinces).Error
	return provinces, err
}

func (r *masterRepository) FindKabupatensByProvID(ctx context.Context, provID int64) ([]entity.Kabupaten, error) {
	var kabs []entity.Kabupaten
	err := r.db.WithContext(ctx).Where("prov_id = ?", provID).Find(&kabs).Error
	return kabs, err
}

func (r *masterRepository) FindKecamatansByKabID(ctx context.Context, kabID int64) ([]entity.Kecamatan, error) {
	var kecs []entity.Kecamatan
	err := r.db.WithContext(ctx).Where("kab_id = ?", kabID).Find(&kecs).Error
	return kecs, err
}

func (r *masterRepository) FindKelurahansByKecID(ctx context.Context, kecID int64) ([]entity.Kelurahan, error) {
	var kels []entity.Kelurahan
	err := r.db.WithContext(ctx).Where("kec_id = ?", kecID).Find(&kels).Error
	return kels, err
}
