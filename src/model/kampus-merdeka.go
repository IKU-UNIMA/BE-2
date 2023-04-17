package model

import "time"

type KampusMerdeka struct {
	ID                      int `gorm:"primaryKey"`
	IdMahasiswa             int
	IdSemester              int
	IdDosenPembimbing       int
	IdKategoriProgram       int
	StatusKeikutsertaan     string `gorm:"type:varchar(255)"`
	KontrakKrs              bool
	JudulAktivitasMahasiswa string `gorm:"type:varchar(255)"`
	SuratTugas              string
	NoSkTugas               string    `gorm:"type:varchar(255)"`
	TanggalSkTugas          time.Time `gorm:"type:date"`
	JenisAnggota            string    `gorm:"type:varchar(255)"`
	Ips                     float32   `gorm:"type:varchar(255)"`
	Ipk                     float32   `gorm:"type:varchar(255)"`
	JumlahSks               int
	TotalSks                int
	BiayaKuliah             float32
	BeritaAcara             string
	Mahasiswa               Mahasiswa         `gorm:"foreignKey:IdMahasiswa"`
	Semester                Semester          `gorm:"foreignKey:IdSemester"`
	DosenPembimbing         Dosen             `gorm:"foreignKey:IdDosenPembimbing"`
	KategoriProgram         KategoriProgramKm `gorm:"foreignKey:IdKategoriProgram"`
}
