package handler

import "time"

// --- Request DTOs ---

type SubmitContractRequest struct {
	MotorID        int64   `json:"motor_id" binding:"required"`
	ProductID      int64   `json:"product_id" binding:"required"`
	NilaiKendaraan float64 `json:"nilai_kendaraan" binding:"required"`
	DpDibayar      float64 `json:"dp_dibayar" binding:"required"`
	TenorBulan     int16   `json:"tenor_bulan" binding:"required"`
}

// --- Response DTOs ---

type ContractResponse struct {
	ContractID      int64      `json:"contract_id"`
	ContractNumber  string     `json:"contract_number"`
	RequestDate     time.Time  `json:"request_date"`
	TanggalAkad     *time.Time `json:"tanggal_akad,omitempty"`
	TenorBulan      int16      `json:"tenor_bulan"`
	NilaiKendaraan  float64    `json:"nilai_kendaraan"`
	DpDibayar       float64    `json:"dp_dibayar"`
	PokokPinjaman   float64    `json:"pokok_pinjaman"`
	TotalPinjaman   float64    `json:"total_pinjaman"`
	CicilanPerBulan float64    `json:"cicilan_per_bulan"`
	Status          string     `json:"status"`
	CustomerID      int64      `json:"customer_id"`
	MotorID         int64      `json:"motor_id"`
	ProductID       int64      `json:"product_id"`
	CreatedAt       time.Time  `json:"created_at"`
}

type TaskProgressResponse struct {
	TaskID          int64      `json:"task_id"`
	TaskName        string     `json:"task_name"`
	Status          string     `json:"status"`
	SequenceNo      int16      `json:"sequence_no"`
	ActualStartdate *time.Time `json:"actual_startdate,omitempty"`
	ActualEnddate   *time.Time `json:"actual_enddate,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type MotorBriefResponse struct {
	MotorID     int64  `json:"motor_id"`
	Merk        string `json:"merk"`
	NomorPolisi string `json:"nomor_polisi"`
}

type MyContractResponse struct {
	ContractID      int64              `json:"contract_id"`
	ContractNumber  string             `json:"contract_number"`
	RequestDate     time.Time          `json:"request_date"`
	Status          string             `json:"status"`
	NilaiKendaraan  float64            `json:"nilai_kendaraan"`
	DpDibayar       float64            `json:"dp_dibayar"`
	CicilanPerBulan float64            `json:"cicilan_per_bulan"`
	TenorBulan      int16              `json:"tenor_bulan"`
	Motor           MotorBriefResponse `json:"motor"`
	CreatedAt       time.Time          `json:"created_at"`
}
