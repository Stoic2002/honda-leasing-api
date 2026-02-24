package entity

import "time"

type PaymentSchedule struct {
	ScheduleID       int64      `gorm:"primaryKey;column:schedule_id"`
	AngsuranKe       int16      `gorm:"column:angsuran_ke"`
	JatuhTempo       time.Time  `gorm:"column:jatuh_tempo"`
	Pokok            float64    `gorm:"column:pokok"`
	Margin           float64    `gorm:"column:margin"`
	TotalTagihan     float64    `gorm:"column:total_tagihan"`
	StatusPembayaran string     `gorm:"column:status_pembayaran"`
	TanggalBayar     *time.Time `gorm:"column:tanggal_bayar"`
	ContractID       int64      `gorm:"column:contract_id"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
}

func (PaymentSchedule) TableName() string {
	return "finance.payment_schedule"
}

type Payment struct {
	PaymentID        int64     `gorm:"primaryKey;column:payment_id"`
	NomorBukti       string    `gorm:"column:nomor_bukti"`
	JumlahBayar      float64   `gorm:"column:jumlah_bayar"`
	TanggalBayar     time.Time `gorm:"column:tanggal_bayar"`
	MetodePembayaran string    `gorm:"column:metode_pembayaran"`
	Provider         *string   `gorm:"column:provider"`
	ContractID       int64     `gorm:"column:contract_id"`
	ScheduleID       *int64    `gorm:"column:schedule_id"`
	CreatedAt        time.Time `gorm:"column:created_at"`
}

func (Payment) TableName() string {
	return "finance.payments"
}
