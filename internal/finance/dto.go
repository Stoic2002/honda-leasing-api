package finance

import "time"

// --- Request DTOs ---

type PaymentRequest struct {
	NomorBukti       string  `json:"nomor_bukti" binding:"required"`
	JumlahBayar      float64 `json:"jumlah_bayar" binding:"required"`
	MetodePembayaran string  `json:"metode_pembayaran" binding:"required"`
	ContractID       int64   `json:"contract_id" binding:"required"`
	ScheduleID       int64   `json:"schedule_id" binding:"required"`
}

// --- Response DTOs ---

type PaymentScheduleResponse struct {
	ScheduleID       int64      `json:"schedule_id"`
	AngsuranKe       int16      `json:"angsuran_ke"`
	JatuhTempo       time.Time  `json:"jatuh_tempo"`
	Pokok            float64    `json:"pokok"`
	Margin           float64    `json:"margin"`
	LateFee          float64    `json:"late_fee"`
	TotalTagihan     float64    `json:"total_tagihan"`
	StatusPembayaran string     `json:"status_pembayaran"`
	TanggalBayar     *time.Time `json:"tanggal_bayar,omitempty"`
}
