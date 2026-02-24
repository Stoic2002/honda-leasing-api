package entity

type Province struct {
	ProvID   int64  `gorm:"primaryKey;column:prov_id"`
	ProvName string `gorm:"column:prov_name"`
}

func (Province) TableName() string {
	return "mst.province"
}

type Kabupaten struct {
	KabID   int64  `gorm:"primaryKey;column:kab_id"`
	KabName string `gorm:"column:kab_name"`
	ProvID  int64  `gorm:"column:prov_id"`
}

func (Kabupaten) TableName() string {
	return "mst.kabupaten"
}

type Kecamatan struct {
	KecID   int64  `gorm:"primaryKey;column:kec_id"`
	KecName string `gorm:"column:kec_name"`
	KabID   int64  `gorm:"column:kab_id"`
}

func (Kecamatan) TableName() string {
	return "mst.kecamatan"
}

type Kelurahan struct {
	KelID   int64  `gorm:"primaryKey;column:kel_id"`
	KelName string `gorm:"column:kel_name"`
	KecID   int64  `gorm:"column:kec_id"`
}

func (Kelurahan) TableName() string {
	return "mst.kelurahan"
}
