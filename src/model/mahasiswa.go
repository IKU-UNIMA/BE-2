package model

type Mahasiswa struct {
	ID           int `gorm:"primaryKey"`
	IdProdi      int
	Nama         string `gorm:"type:varchar(255)"`
	Nim          string `gorm:"type:varchar(255);unique"`
	JenisKelamin string `gorm:"type:varchar(1)"`
	Prodi        Prodi  `gorm:"foreignKey:IdProdi;constraint:OnDelete:CASCADE"`
}
