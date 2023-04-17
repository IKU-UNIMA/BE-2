package response

type (
	Prestasi struct {
		ID          int                `json:"id"`
		IdMahasiswa int                `json:"-"`
		IdSemester  int                `json:"-"`
		Nama        string             `json:"nama"`
		Mahasiswa   MahasiswaReference `gorm:"foreignKey:IdMahasiswa;constraint:OnDelete:CASCADE" json:"mahasiswa"`
		Semester    Semester           `gorm:"foreignKey:IdSemester;constraint:OnDelete:CASCADE" json:"semester"`
	}

	DetailPrestasi struct {
		ID                int                `json:"id"`
		IdMahasiswa       int                `json:"-"`
		IdSemester        int                `json:"-"`
		IdDosenPembimbing int                `json:"-"`
		Nama              string             `json:"nama"`
		TingkatPrestasi   string             `json:"tingkat_prestasi"`
		Penyelenggara     string             `json:"penyelenggara"`
		Peringkat         string             `json:"peringkat"`
		Sertifikat        string             `json:"sertifikat"`
		Mahasiswa         MahasiswaReference `gorm:"foreignKey:IdMahasiswa" json:"mahasiswa"`
		Semester          Semester           `gorm:"foreignKey:IdSemester" json:"semester"`
		DosenPembimbing   Dosen              `gorm:"foreignKey:IdDosenPembimbing" json:"dosen_pembimbing"`
	}
)
