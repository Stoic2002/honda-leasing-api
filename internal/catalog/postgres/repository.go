package postgres

import (
	"context"
	"fmt"

	"honda-leasing-api/internal/domain/contract"
	"honda-leasing-api/internal/domain/entity"

	"gorm.io/gorm"
)

type catalogRepository struct {
	db *gorm.DB
}

func NewCatalogRepository(db *gorm.DB) contract.CatalogRepository {
	return &catalogRepository{db: db}
}

func (r *catalogRepository) FindAllMotors(ctx context.Context, filter contract.CatalogFilter, pagination contract.PaginationFilter) ([]entity.Motor, int64, error) {
	var motors []entity.Motor
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Motor{}).
		Preload("MotorType").
		Preload("Assets")

	if filter.Search != "" {
		searchTerm := fmt.Sprintf("%%%s%%", filter.Search)
		query = query.Where("merk ILIKE ? OR motor_type_name ILIKE ?", searchTerm, searchTerm)
		// Note we can't search MotorTypeName directly without join, modifying query to handle Type join
	}

	if filter.MotorType != "" && filter.MotorType != "all" {
		query = query.Joins("JOIN dealer.motor_types mt ON mt.moty_id = dealer.motors.motor_moty_id").
			Where("mt.moty_name ILIKE ?", filter.MotorType)
	}

	// Count total records for pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.Limit
	query = query.Offset(offset).Limit(pagination.Limit).Order("motor_id ASC")

	if err := query.Find(&motors).Error; err != nil {
		return nil, 0, err
	}

	return motors, total, nil
}

func (r *catalogRepository) FindMotorByID(ctx context.Context, id int64) (*entity.Motor, error) {
	var motor entity.Motor
	if err := r.db.WithContext(ctx).
		Preload("MotorType").
		Preload("Assets").
		First(&motor, id).Error; err != nil {
		return nil, err
	}
	return &motor, nil
}

func (r *catalogRepository) GetLeasingProducts(ctx context.Context) ([]entity.LeasingProduct, error) {
	var products []entity.LeasingProduct
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
