package response

type Operator struct {
	ID      int            `json:"id"`
	IdProdi int            `json:"-"`
	Nama    string         `json:"nama"`
	Nip     string         `json:"nip"`
	Email   string         `json:"email"`
	Prodi   ProdiReference `gorm:"foreignKey:IdProdi" json:"prodi"`
}
