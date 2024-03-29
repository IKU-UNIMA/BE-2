package response

type (
	Mahasiswa struct {
		ID           int    `json:"id"`
		IdProdi      int    `json:"-"`
		Nama         string `json:"nama"`
		Nim          string `json:"nim"`
		Email        string `json:"email"`
		JenisKelamin string `json:"jenis_kelamin"`
		Prodi        Prodi  `gorm:"foreignKey:IdProdi" json:"prodi"`
	}

	MahasiswaReference struct {
		ID      int    `json:"id"`
		IdProdi int    `json:"-"`
		Nama    string `json:"nama"`
		Nim     string `json:"nim"`
		Prodi   Prodi  `gorm:"foreignKey:IdProdi" json:"prodi"`
	}
)

func (MahasiswaReference) TableName() string {
	return "mahasiswa"
}
