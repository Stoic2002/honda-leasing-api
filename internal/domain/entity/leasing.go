package entity

import "time"

type LeasingProduct struct {
	ProductID   int64     `gorm:"primaryKey;column:product_id"`
	KodeProduk  string    `gorm:"column:kode_produk"`
	NamaProduk  string    `gorm:"column:nama_produk"`
	TenorBulan  int16     `gorm:"column:tenor_bulan"`
	DpPersenMin float64   `gorm:"column:dp_persen_min"`
	DpPersenMax float64   `gorm:"column:dp_persen_max"`
	BungaFlat   float64   `gorm:"column:bunga_flat"`
	AdminFee    float64   `gorm:"column:admin_fee"`
	Asuransi    bool      `gorm:"column:asuransi"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (LeasingProduct) TableName() string {
	return "leasing.leasing_product"
}

type LeasingContract struct {
	ContractID        int64      `gorm:"primaryKey;column:contract_id"`
	ContractNumber    string     `gorm:"column:contract_number"`
	RequestDate       time.Time  `gorm:"column:request_date"`
	TanggalAkad       *time.Time `gorm:"column:tanggal_akad"`
	TanggalMulaiCicil *time.Time `gorm:"column:tanggal_mulai_cicil"`
	TenorBulan        int16      `gorm:"column:tenor_bulan"`
	NilaiKendaraan    float64    `gorm:"column:nilai_kendaraan"`
	DpDibayar         float64    `gorm:"column:dp_dibayar"`
	PokokPinjaman     float64    `gorm:"column:pokok_pinjaman"`
	TotalPinjaman     float64    `gorm:"column:total_pinjaman"`
	CicilanPerBulan   float64    `gorm:"column:cicilan_per_bulan"`
	Status            string     `gorm:"column:status"`
	CustomerID        int64      `gorm:"column:customer_id"`
	MotorID           int64      `gorm:"column:motor_id"`
	ProductID         int64      `gorm:"column:product_id"`
	CreatedAt         time.Time  `gorm:"column:created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at"`

	// Relationships
	Customer Customer       `gorm:"foreignKey:CustomerID;references:CustomerID"`
	Motor    Motor          `gorm:"foreignKey:MotorID;references:MotorID"`
	Product  LeasingProduct `gorm:"foreignKey:ProductID;references:ProductID"`
}

func (LeasingContract) TableName() string {
	return "leasing.leasing_contract"
}
