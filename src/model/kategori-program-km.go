package model

type KategoriProgramKm struct {
	ID   int    `gorm:"foreignKey"`
	Nama string `gorm:"type:varchar(255)"`
}
