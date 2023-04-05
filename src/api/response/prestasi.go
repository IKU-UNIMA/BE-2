package response

type (
	Prestasi struct {
		ID          int                `json:"id"`
		IdMahasiswa int                `json:"-"`
		IdProdi     int                `json:"-"`
		IdSemester  int                `json:"-"`
		Nama        string             `json:"nama"`
		Mahasiswa   MahasiswaReference `gorm:"foreignKey:IdMahasiswa;constraint:OnDelete:CASCADE"`
		Prodi       ProdiReference     `gorm:"foreignKey:IdProdi;constraint:OnDelete: SET NULL"`
		Semester    Semester           `gorm:"foreignKey:IdSemester;constraint:OnDelete:CASCADE"`
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
		Mahasiswa         MahasiswaReference `gorm:"foreignKey:IdMahasiswa"`
		Fakultas          Fakultas           `gorm:"foreignKey:IdFakultas"`
		Prodi             ProdiReference     `gorm:"foreignKey:IdProdi"`
		Semester          Semester           `gorm:"foreignKey:IdSemester"`
		DosenPembimbing   Dosen              `gorm:"foreignKey:IdDosenPembimbing"`
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
