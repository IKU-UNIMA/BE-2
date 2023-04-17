package response

type KampusMerdeka struct {
	ID                      int                `json:"id"`
	IdMahasiswa             int                `json:"-"`
	IdSemester              int                `json:"-"`
	IdKategoriProgram       int                `json:"-"`
	Mahasiswa               MahasiswaReference `gorm:"foreignKey:IdMahasiswa" json:"mahasiswa"`
	Semester                Semester           `gorm:"foreignKey:IdSemester" json:"semester"`
	JudulAktivitasMahasiswa string             `json:"judul_aktivitas_mahasiswa"`
	KategoriProgram         KategoriProgramKm  `gorm:"foreignKey:IdKategoriProgram" json:"kategori_program"`
}

type DetailKampusMerdeka struct {
	ID                      int                `json:"id"`
	IdMahasiswa             int                `json:"-"`
	IdSemester              int                `json:"-"`
	IdKategoriProgram       int                `json:"-"`
	IdDosenPembimbing       int                `json:"-"`
	Mahasiswa               MahasiswaReference `gorm:"foreignKey:IdMahasiswa" json:"mahasiswa"`
	Semester                Semester           `gorm:"foreignKey:IdSemester" json:"semester"`
	JudulAktivitasMahasiswa string             `json:"judul_aktivitas_mahasiswa"`
	KategoriProgram         KategoriProgramKm  `gorm:"foreignKey:IdKategoriProgram" json:"kategori_program"`
	StatusKeikutsertaan     string             `json:"status_keikutsertaan"`
	JenisAnggota            string             `json:"jenis_anggota"`
	DosenPembimbing         Dosen              `gorm:"foreignKey:IdDosenPembimbing" json:"dosen_pembimbing"`
	NoSkTugas               string             `json:"no_sk_tugas"`
	TanggalSkTugas          string             `json:"tanggal_sk_tugas"`
	SuratTugas              string             `json:"surat_tugas"`
	KontrakKrs              bool               `json:"kontrak_krs"`
	Ips                     float32            `json:"ips"`
	Ipk                     float32            `json:"Ipk"`
	JumlahSks               int                `json:"jumlah_sks"`
	TotalSks                int                `json:"total_sks"`
	BiayaKuliah             float32            `json:"biaya_kuliah"`
	BeritaAcara             string             `json:"berita_acara"`
}
