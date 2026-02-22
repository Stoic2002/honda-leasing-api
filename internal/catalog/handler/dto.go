package handler

import "time"

// --- Motor DTOs ---

type MotorTypeResponse struct {
	MotyID   int64  `json:"moty_id"`
	MotyName string `json:"moty_name"`
}

type MotorAssetResponse struct {
	MoasID   int64   `json:"moas_id"`
	FileName string  `json:"file_name"`
	FileSize float64 `json:"file_size"`
	FileType string  `json:"file_type"`
	FileURL  string  `json:"file_url"`
}

type MotorResponse struct {
	MotorID     int64                `json:"motor_id"`
	Merk        string               `json:"merk"`
	Tahun       int16                `json:"tahun"`
	Warna       string               `json:"warna"`
	NomorRangka string               `json:"nomor_rangka"`
	NomorMesin  string               `json:"nomor_mesin"`
	CCMesin     string               `json:"cc_mesin"`
	NomorPolisi string               `json:"nomor_polisi"`
	StatusUnit  string               `json:"status_unit"`
	HargaOTR    float64              `json:"harga_otr"`
	MotorType   MotorTypeResponse    `json:"motor_type"`
	Assets      []MotorAssetResponse `json:"assets"`
	CreatedAt   time.Time            `json:"created_at"`
}

// --- Leasing Product DTOs ---

type LeasingProductResponse struct {
	ProductID   int64   `json:"product_id"`
	KodeProduk  string  `json:"kode_produk"`
	NamaProduk  string  `json:"nama_produk"`
	TenorBulan  int16   `json:"tenor_bulan"`
	DpPersenMin float64 `json:"dp_persen_min"`
	DpPersenMax float64 `json:"dp_persen_max"`
	BungaFlat   float64 `json:"bunga_flat"`
	AdminFee    float64 `json:"admin_fee"`
	Asuransi    bool    `json:"asuransi"`
}
