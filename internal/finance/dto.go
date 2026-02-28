package finance

import "time"

// --- Request DTOs ---

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
