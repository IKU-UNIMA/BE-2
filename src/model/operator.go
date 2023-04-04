package model

type Operator struct {
	ID      int `gorm:"primaryKey"`
	IdProdi int
	Nama    string `gorm:"type:varchar(255)"`
	Nip     string `gorm:"type:varchar(255);unique"`
	Prodi   Prodi  `gorm:"foreignKey:IdProdi;constraint:OnDelete:CASCADE"`
}
