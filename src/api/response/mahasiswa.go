package response

type Mahasiswa struct {
	ID           int            `json:"id"`
	IdProdi      int            `json:"-"`
	Nama         string         `json:"nama"`
	Nim          string         `json:"nim"`
	Email        string         `json:"email"`
	JenisKelamin string         `json:"jenis_kelamin"`
	Prodi        ProdiReference `gorm:"foreignKey:IdProdi" json:"prodi"`
}
