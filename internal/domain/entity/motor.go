package entity

import "time"

type Motor struct {
	MotorID     int64     `gorm:"primaryKey;column:motor_id"`
	Merk        string    `gorm:"column:merk"`
	Tahun       int16     `gorm:"column:tahun"`
	Warna       string    `gorm:"column:warna"`
	NomorRangka string    `gorm:"column:nomor_rangka"`
	NomorMesin  string    `gorm:"column:nomor_mesin"`
	CCMesin     string    `gorm:"column:cc_mesin"`
	NomorPolisi string    `gorm:"column:nomor_polisi"`
	StatusUnit  string    `gorm:"column:status_unit"`
	HargaOTR    float64   `gorm:"column:harga_otr"`
	MotorMotyID int64     `gorm:"column:motor_moty_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`

	// Relationships
	MotorType MotorType    `gorm:"foreignKey:MotorMotyID;references:MotyID"`
	Assets    []MotorAsset `gorm:"foreignKey:MoasMotorID;references:MotorID"`
}

func (Motor) TableName() string {
	return "dealer.motors"
}

type MotorType struct {
	MotyID   int64  `gorm:"primaryKey;column:moty_id"`
	MotyName string `gorm:"column:moty_name"`
}

func (MotorType) TableName() string {
	return "dealer.motor_types"
}

type MotorAsset struct {
	MoasID      int64     `gorm:"primaryKey;column:moas_id"`
	FileName    string    `gorm:"column:file_name"`
	FileSize    float64   `gorm:"column:file_size"`
	FileType    string    `gorm:"column:file_type"`
	FileURL     string    `gorm:"column:file_url"`
	MoasMotorID int64     `gorm:"column:moas_motor_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (MotorAsset) TableName() string {
	return "dealer.motor_assets"
}
