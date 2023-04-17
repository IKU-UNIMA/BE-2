package model

type Prestasi struct {
	ID                int `gorm:"primaryKey"`
	IdMahasiswa       int
	IdSemester        int
	IdDosenPembimbing int       `gorm:"default:null"`
	Nama              string    `gorm:"type:varchar(255)"`
	TingkatPrestasi   string    `gorm:"type:varchar(60)"`
	Penyelenggara     string    `gorm:"type:text"`
	Peringkat         string    `gorm:"type:varchar(30)"`
	Sertifikat        string    `gorm:"type:text"`
	Mahasiswa         Mahasiswa `gorm:"foreignKey:IdMahasiswa;constraint:OnDelete:CASCADE"`
	Semester          Semester  `gorm:"foreignKey:IdSemester;constraint:OnDelete:CASCADE"`
	DosenPembimbing   Dosen     `gorm:"foreignKey:IdDosenPembimbing;constraint:OnDelete:SET NULL"`
}
