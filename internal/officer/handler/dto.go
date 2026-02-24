package handler

import "time"

// --- Response DTOs ---

type CustomerBriefResponse struct {
	CustomerID  int64  `json:"customer_id"`
	NamaLengkap string `json:"nama_lengkap"`
	NoHp        string `json:"no_hp"`
	Email       string `json:"email"`
}

type MotorBriefResponse struct {
	MotorID     int64  `json:"motor_id"`
	Merk        string `json:"merk"`
	NomorPolisi string `json:"nomor_polisi"`
}

type IncomingOrderResponse struct {
	ContractID     int64                 `json:"contract_id"`
	ContractNumber string                `json:"contract_number"`
	RequestDate    time.Time             `json:"request_date"`
	Status         string                `json:"status"`
	NilaiKendaraan float64               `json:"nilai_kendaraan"`
	DpDibayar      float64               `json:"dp_dibayar"`
	Customer       CustomerBriefResponse `json:"customer"`
	Motor          MotorBriefResponse    `json:"motor"`
	CreatedAt      time.Time             `json:"created_at"`
}

type OfficerTaskResponse struct {
	TaskID          int64      `json:"task_id"`
	TaskName        string     `json:"task_name"`
	Status          string     `json:"status"`
	SequenceNo      int16      `json:"sequence_no"`
	ContractID      int64      `json:"contract_id"`
	ActualStartdate *time.Time `json:"actual_startdate,omitempty"`
	ActualEnddate   *time.Time `json:"actual_enddate,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}
