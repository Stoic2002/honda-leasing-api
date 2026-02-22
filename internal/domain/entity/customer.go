package entity

import "time"

type Customer struct {
	CustomerID   int64      `gorm:"primaryKey;column:customer_id"`
	UserID       *int64     `gorm:"column:user_id"`
	Nik          string     `gorm:"column:nik"`
	NamaLengkap  string     `gorm:"column:nama_lengkap"`
	TanggalLahir *time.Time `gorm:"column:tanggal_lahir"`
	NoHp         string     `gorm:"column:no_hp"`
	Email        string     `gorm:"column:email"`
	Pekerjaan    string     `gorm:"column:pekerjaan"`
	Perusahaan   string     `gorm:"column:perusahaan"`
	Salary       float64    `gorm:"column:salary"`
	LocationID   *int64     `gorm:"column:location_id"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
}

func (Customer) TableName() string {
	return "dealer.customers"
}
