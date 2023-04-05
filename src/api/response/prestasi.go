package response

type (
	Prestasi struct {
		ID          int                `json:"id"`
		IdMahasiswa int                `json:"-"`
		IdProdi     int                `json:"-"`
		IdSemester  int                `json:"-"`
		Nama        string             `json:"nama"`
		Mahasiswa   MahasiswaReference `gorm:"foreignKey:IdMahasiswa;constraint:OnDelete:CASCADE" json:"mahasiswa"`
		Prodi       ProdiReference     `gorm:"foreignKey:IdProdi;constraint:OnDelete: SET NULL" json:"prodi"`
		Semester    Semester           `gorm:"foreignKey:IdSemester;constraint:OnDelete:CASCADE" json:"semester"`
	}

	DetailPrestasi struct {
		ID                int                `json:"id"`
		IdMahasiswa       int                `json:"-"`
		IdFakultas        int                `json:"-"`
		IdProdi           int                `json:"-"`
		IdSemester        int                `json:"-"`
		IdDosenPembimbing int                `json:"-"`
		Nama              string             `json:"nama"`
		TingkatPrestasi   string             `json:"tingkat_prestasi"`
		Penyelenggara     string             `json:"penyelenggara"`
		Peringkat         string             `json:"peringkat"`
		Sertifikat        string             `json:"sertifikat"`
		Mahasiswa         MahasiswaReference `gorm:"foreignKey:IdMahasiswa" json:"mahasiswa"`
		Fakultas          Fakultas           `gorm:"foreignKey:IdFakultas" json:"fakultas"`
		Prodi             ProdiReference     `gorm:"foreignKey:IdProdi" json:"prodi"`
		Semester          Semester           `gorm:"foreignKey:IdSemester" json:"semester"`
		DosenPembimbing   Dosen              `gorm:"foreignKey:IdDosenPembimbing" json:"dosen_pembimbing"`
	}

	MahasiswaReference struct {
		ID   int    `json:"id"`
		Nama string `json:"nama"`
		Nim  string `json:"nim"`
	}
)

func (MahasiswaReference) TableName() string {
	return "mahasiswa"
}
